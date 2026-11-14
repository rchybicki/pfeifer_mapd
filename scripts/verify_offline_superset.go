//go:build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/encoding/text"
	"capnproto.org/go/capnp/v3/schemas"
	"pfeifer.dev/mapd/cereal/offline"
)

func waysBySignature(path string) (map[string]int, int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, err
	}

	msg, err := capnp.UnmarshalPacked(data)
	if err != nil {
		return nil, 0, err
	}
	defer msg.Release()

	tile, err := offline.ReadRootOffline(msg)
	if err != nil {
		return nil, 0, err
	}
	ways, err := tile.Ways()
	if err != nil {
		return nil, 0, err
	}

	signatures := make(map[string]int, ways.Len())
	for i := range ways.Len() {
		signature, err := text.Marshal(offline.Way_TypeID, capnp.Struct(ways.At(i)))
		if err != nil {
			return nil, 0, err
		}
		signatures[signature]++
	}
	return signatures, ways.Len(), nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: go run verify_offline_superset.go OLD_DIR NEW_DIR")
		os.Exit(2)
	}

	offline.RegisterSchema(schemas.DefaultRegistry)

	oldDir := os.Args[1]
	newDir := os.Args[2]
	oldWays := 0
	newWays := 0
	checkedTiles := 0

	err := filepath.WalkDir(newDir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}

		relativePath, err := filepath.Rel(newDir, path)
		if err != nil {
			return err
		}
		oldPath := filepath.Join(oldDir, relativePath)

		oldSignatures, oldCount, err := waysBySignature(oldPath)
		if err != nil {
			return fmt.Errorf("read old tile %s: %w", relativePath, err)
		}
		newSignatures, newCount, err := waysBySignature(path)
		if err != nil {
			return fmt.Errorf("read new tile %s: %w", relativePath, err)
		}

		for signature, oldOccurrences := range oldSignatures {
			if newSignatures[signature] < oldOccurrences {
				return fmt.Errorf("tile %s lost %d occurrence(s) of an existing way", relativePath, oldOccurrences-newSignatures[signature])
			}
		}

		checkedTiles++
		oldWays += oldCount
		newWays += newCount
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("verified %d tiles: all %d existing ways retained; %d ways after addition\n", checkedTiles, oldWays, newWays)
}
