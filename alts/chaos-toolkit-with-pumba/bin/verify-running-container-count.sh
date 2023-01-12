#!/usr/bin/env bash

count="$(docker ps -f "name=$1" --format '{{.ID}}' | wc -l)"
if [ "$count" -ne $2 ]; then
  echo "Expected $2 container(s), got $count" >&2
  exit 1
fi
