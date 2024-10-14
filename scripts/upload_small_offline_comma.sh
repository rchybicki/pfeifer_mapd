#!/bin/bash

rsync -av --exclude='*.tar.gz' --exclude='._*' offline/ commawifi:/data/media/0/osm/offline/

