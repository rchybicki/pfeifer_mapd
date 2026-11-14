# Mapd v2.x
This branch contains the docs and code for v2.x of mapd. v2.x is incompatible
with v1.x. v1.x can be found [here](https://github.com/pfeiferj/mapd/tree/v1.x).

Mapd is a component for use in forks of comma.ai's
[openpilot](https://github.com/commaai/openpilot). It uses openstreetmap data to
provide speed limit and curve speed adjustments to the openpilot fork.

For the current FrogPilot testing integration, this repo is used to build the
local `mapd` generator, generate offline map data, and upload offline tiles.
The runtime `mapd` binary that FrogPilot launches is tracked in the sibling
`openpilot-rch` repo at `frogpilot/navigation/mapd`, so runtime binary updates
should be made and deployed from that repo instead of pushing a binary directly
from here.

## Route Speed Overrides
Route speed overrides are defined in [maps/speed_overrides.go](./maps/speed_overrides.go)
and are applied when offline map tiles are generated. Editing those values does
not change runtime behavior until the offline tiles are regenerated from the
same OSM PBF extract.

## Custom Europe Trip Maps
Use `scripts/download_extracts.sh poland germany czech-republic austria italy`
to download the Geofabrik extracts locally without touching the device. The
script skips existing extracts unless `MAPD_FORCE_DOWNLOAD=1` is set.

Use `scripts/generate_offline_Poland_Germany_Czech_Austria_Italy.sh` to merge
the Poland, Germany, Czech Republic, Austria, and Italy extracts, filter them,
and generate a combined offline tile set under `scripts/offline`. It covers
longitude `5..25` and latitude `35..55`, which contains Poland, Germany, the
Czech Republic, Austria, and Italy. It does not upload by default; set
`MAPD_UPLOAD=1` only when the generated local tiles should be pushed to the
comma.

`scripts/download_generate_Poland_Germany_Czech_Austria_Italy.sh` runs both
steps and refreshes all extracts by default so border objects come from the same
Geofabrik snapshot. Generation still requires `osmium-tool`, matching the
existing map generation scripts.

The older `scripts/generate_offline_Poland_Germany.sh` and
`scripts/download_generate_Poland_Germany.sh` scripts are still available for a
smaller Poland and Germany-only tile set.

### Selective Slovenia addition
Use `scripts/download_generate_Slovenia.sh` to download only the Slovenia
Geofabrik extract and generate only the quarter-degree tiles enclosing Slovenia.
The generator combines Slovenia with the existing five-country merged source so
Austria and Italy border data is retained, and it applies the same route speed
overrides as the full generator.

The output is written to `scripts/offline-slovenia` and is not uploaded by
default. Set `MAPD_UPLOAD=1` to call
`scripts/upload_selective_offline_comma.sh`, which adds only the generated files
and deliberately does not use `rsync --delete`. This path does not rebuild or
deploy the device's runtime `mapd` binary.

Before uploading, `go run scripts/verify_offline_superset.go
scripts/offline scripts/offline-slovenia` verifies that every complete road
record in the previous tiles is still present in the Slovenia-enhanced tiles.

### Selective Poland refresh
Use `scripts/update_generate_Poland.sh` to update only the cached Poland extract
through its Geofabrik replication feed and then generate the 984 quarter-degree
tiles enclosing Poland. The generator combines the updated Poland source with
the unchanged Germany, Czech Republic, and Austria extracts so border roads are
retained and applies the same route speed overrides as the full generator.

The output is written to `scripts/offline-poland-update` and is not uploaded by
default. Set `MAPD_UPLOAD=1` to use the checksum-aware selective uploader. The
uploader compresses rsync traffic, does not use `--delete`, and does not rebuild
or deploy the runtime `mapd` binary.

## openpilot fork integration
openpilot fork integration is described in
[docs/integration.md](./docs/integration.md).

## Disclaimers and Acknowledgements
* "openpilot" and "msgq" are trademarks of comma.ai. mapd is in no way affiliated with
  comma.ai, openpilot, or msgq.

* The capnp definitions in the cereal package are modifications of the
  definitions in the cereal module in
  [openpilot](https://github.com/commaai/openpilot/tree/master/cereal). openpilot
  is released under the MIT license.

* Many of the ideas in this code were built by referencing the
  [move-fast implementation](https://github.com/move-fast/openpilot/tree/b170d1bc123a0cf2b872050fbd5e2eecd1b23e22/selfdrive/mapd/lib)
  of an openpilot map daemon. The move-fast implementation was released under the
  MIT license.

* Other ideas in this implementation came from the wider openpilot fork
  community. Some notable mentions are:
    * [sunnypilot](https://github.com/sunnypilot/sunnypilot)
    * [frogpilot](https://github.com/FrogAi/FrogPilot)
    * [dragonpilot](https://github.com/dragonpilot/dragonpilot)
