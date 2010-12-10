#!/usr/bin/env bash

echo "####### Building Libs #######"
cd lib
./build-libs
echo "###### Building Filters #####"
cd ../filters
./build-filters
echo "####### Building Main #######"
cd ../main
make	
echo "##### Building Manpages #####"
cd ../doc
make