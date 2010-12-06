#!/usr/bin/env bash

echo "####### Building Libs #######"
cd lib
./build-libs
echo "####### Building Filters ####"
cd ../filters-src
./build-filters
echo "####### Building Main ######"
cd ../main
make	
cd ..

chmod -R g+rwx $GOROOT/pkg/linux_amd64/afp* 2> /dev/null
