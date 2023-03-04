#!/bin/bash

MAPS="de_overpass de_mirage de_vertigo de_vertigo_lower de_nuke de_nuke_lower de_cache de_inferno de_train de_dust2 de_ancient de_anubis"

mkdir -p overviews
for MAP in ${MAPS}; do
   echo "downloading ${MAP} ..."
   curl -s -o "overviews/${MAP}.jpg" "https://raw.githubusercontent.com/zoidbergwill/csgo-overviews/master/overviews/${MAP}.jpg"
done
