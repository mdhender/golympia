#!/bin/bash

if [ ! -f "item" ]; then
  echo "error: missing 'item' file"
  exit 2
fi

awk -f items.awk < item > items.json

exit $?
