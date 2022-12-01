#!/usr/bin/env bash

# NOTE: the network interfaces seem to be assigned alphabetically

# show eth0 (control) network interface using ip command
ip addr show eth0

go run main.go
