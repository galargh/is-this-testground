#!/usr/bin/env bash

# NOTE: the network interfaces seem to be assigned alphabetically

# show eth0 (control) network interface using ip command
ip addr show eth0
# show eth1 (data) network interface using ip command
ip addr show eth1

# delay packets for 100ms
# https://www.excentis.com/blog/use-linux-traffic-control-as-impairment-node-in-a-test-environment-part-2/
tc qdisc add dev eth1 root netem delay 1000ms

go run main.go
