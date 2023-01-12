#!/usr/bin/env bash

docker_ids="$(docker ps -f "name=$1" --format '{{.ID}}')"
ipfs_ids=""
while read -r docker_id; do
  ipfs_id="$(docker exec "$docker_id" ipfs id -f '<id>')"
  while read -r other_docker_id; do
    if [ "$docker_id" = "$other_docker_id" ]; then
      continue
    fi
    ipfs_ids="$(docker exec "$other_docker_id" ipfs swarm peers)"
    if ! echo "$ipfs_ids" | grep -q "$ipfs_id"; then
      echo "Expected $ipfs_id to be connected to $other_docker_id" >&2
      exit 1
    fi
  done <<< "$docker_ids"
done <<< "$docker_ids"
