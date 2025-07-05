#!/bin/bash

wget -O europe.osm.pbf https://download.bbbike.org/osm/planet/sub-planet-daily/europe.osm.pbf


./generate_and_update_Poland.sh
