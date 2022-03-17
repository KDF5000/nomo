#!/bin/bash

DIR=$(cd "$(dirname "$0")"; pwd)
cd $DIR/bin; NOMO_LOG_FILE=$DIR/nomo.log ./nomo