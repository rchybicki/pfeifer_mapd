#!/bin/bash
MIN_LON=14
MIN_LAT=49
MAX_LON=24
MAX_LAT=54

./filter_Europe.sh
rm -r offline

for (( i = $MIN_LON; i < $MAX_LON; i += 20 )) 
do
  for (( j = $MIN_LAT; j < $MAX_LAT; j += 20 )) 
  do
    max_lon=$(($i+20))
    max_lat=$(($j+20))
    echo "$i $j $max_lon $max_lat"
    ./extract_box.sh $i $j $max_lon $max_lat
    ./add_locations.sh
    ./mapd --generate --minlat $j --minlon $i --maxlat $max_lat --maxlon $max_lon
    #./compress_offline.sh
    ./upload_small_offline_comma.sh
    rm -r offline
  done

done
