#!/bin/bash
for i in $(seq 0 14); do
  ~/pastelnetwork/pastel/src/pastel-cli -datadir=/home/nd/pastelnetwork/mock-networks/regtest/unzipped/node${i} getinfo | grep connections
done
