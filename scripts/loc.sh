#!/bin/bash

if [ ! -f "loc" ]; then
  echo "error: missing 'loc' file"
  exit 2
fi

awk -f loc.awk < loc > loc.json

exit $?
