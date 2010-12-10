#!/usr/bin/env bash

#Copyright (c) 2010 Go Fightclub Authors
#This source code is released under the terms of the
#MIT license. Please see the file LICENSE for license details.

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