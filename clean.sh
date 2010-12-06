#!/usr/bin/env bash

echo "####### Building Libs #######"
cd lib
./clean-libs
echo "###### Building Filters #####"
cd ../filters
./clean-filters
echo "####### Building Main #######"
cd ../main
make clean	
echo "##### Building Manpages #####"
cd ../doc
make clean