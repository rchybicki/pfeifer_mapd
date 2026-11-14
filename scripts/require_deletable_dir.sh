#!/bin/bash
# Sourced by the offline generators before they rm -rf a directory.

# Absolute path with every existing component resolved through symlinks.
# Components that do not exist yet are appended unchanged.
physical_path() {
  local existing="$1" missing=""
  while [[ "$existing" == ?*/ ]]; do
    existing="${existing%/}"
  done
  [[ "$existing" == /* ]] || existing="$PWD/$existing"
  while [[ ! -d "$existing" ]]; do
    missing="/${existing##*/}$missing"
    existing="${existing%/*}"
    existing="${existing:-/}"
  done
  existing="$(cd -- "$existing" && pwd -P)" || return
  printf '%s%s\n' "${existing%/}" "$missing"
}

# require_deletable_dir NAME VALUE [full]
# Exits unless what rm -rf VALUE removes lies strictly inside this scripts
# directory and clear of the full host tile set, scripts/offline. A full
# generator passes "full" to replace that set as a whole.
require_deletable_dir() {
  local name="$1" value="$2" scripts entry tiles target
  scripts="$(physical_path "$(dirname -- "${BASH_SOURCE[0]}")")"
  entry="$scripts/offline"
  tiles="$(physical_path "$entry")"
  # rm -rf removes a symlink named without a trailing slash, not its target.
  if [[ "$value" != */ && -L "$value" ]]; then
    target="$(physical_path "$(dirname -- "$value")")/${value##*/}"
  else
    target="$(physical_path "$value")"
  fi
  case "/$value/" in
    */../*) ;;
    *)
      if [[ "$target" == "$entry" ]]; then
        [[ "${3:-}" == full ]] && return 0
      elif [[ "$target" == "$scripts"/* && "$target/" != "$tiles"/* && "$tiles/" != "$target"/* ]]; then
        return 0
      fi
      ;;
  esac
  echo "refusing unsafe $name: $value (removes $target)" >&2
  exit 1
}
