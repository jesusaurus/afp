#!/usr/bin/env bash

#Copyright (c) 2010 AFP Authors
#This source code is released under the terms of the
#MIT license. Please see the file LICENSE for license details.

cd lib
./clean-libs
cd ../filters
./clean-filters
cd ../main
make clean	
cd ..