# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

### Environment Setup
This project uses [earthly](https://github.com/earthly/earthly/) for its build system. To install earthly follow the instructions at the [get earthly page](https://earthly.dev/get-earthly).

### Common Commands
- **Format code**: `earthly +format`
- **Lint code**: `earthly +lint`
- **Run tests**: `earthly +test`
- **Update snapshot tests**: `earthly +update-snapshots`
- **Build capnp files**: `earthly +compile-capnp`
- **Build release binary** (for ARM64/comma device): `earthly +build-release`
- **Build binary** (for current architecture): `earthly +build`

### Script Usage
- To update the OpenPilot binary:
  ```bash
  earthly +build
  cd scripts
  ./upload_mapd_comma.sh
  ```
- To update maps:
  ```bash
  ./scripts/generate_and_update_Poland.sh
  ```

## Architecture Overview

### Core Components

1. **mapd** - The main daemon that provides OpenPilot with map data
   - Provides speed limits, curvature data, road names, and more to the OpenPilot system
   - Communicates via memory params in `/dev/shm/params`

2. **Data Flow**:
   - Reads GPS position data from OpenPilot
   - Loads map data for the current location
   - Processes route information (current way, next ways)
   - Outputs structured data for OpenPilot to consume

3. **Map Data Management**:
   - Downloads map data from OpenStreetMap (OSM)
   - Supports downloading by bounding box or pre-defined locations
   - Processes the data into an optimized offline format

### Key Features

- **Turn Speed Control**: Calculates target velocities based on road curvature data
- **Speed Limit Detection**: Determines speed limits based on OSM data, including directional limits
- **Hazard Detection**: Provides information about hazards marked in OSM
- **Dynamic Configuration**: Configurable log levels and parameters at runtime

## Working with This Codebase

### Key Files
- `mapd.go`: Main loop and primary logic
- `download.go`: Map data downloading functionality
- `way.go`: Road/way data structures and operations
- `math.go`: Mathematical computations for curvature, etc.
- `params.go`: Parameter handling for communication with OpenPilot
- `offline.capnp`: Schema for the offline map data format

### Testing
- Tests use snapshot-based assertions
- Math-related tests are in `math_test.go`

### Inputs and Outputs
- Inputs: GPS coordinates and bearing from `/dev/shm/params/LastGPSPosition`
- Outputs: Various data like speed limits, road names, and curvatures written to `/dev/shm/params/`
- Details on inputs/outputs are in the docs directory (`docs/inputs.md` and `docs/outputs.md`)