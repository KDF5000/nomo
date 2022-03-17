#!/bin/bash
# used for ansible

WORKDIR=/opt/openhex/nomo
export NOMO_LOG_FILE=${WORKDIR}/logs/nomo.log 
exec ${WORKDIR}/bin/nomo > std.log  2>&1
