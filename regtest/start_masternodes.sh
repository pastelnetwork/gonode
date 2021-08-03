#!/bin/sh
for i in $(seq 0 9); do
  ~/pastelnetwork/pastel/src/pastel-cli -datadir=/home/nd/pastelnetwork/mock-networks/regtest/unzipped/node14 masternode start-alias mn${i}
done

# echo "generate 100 blocks"
# ~/pastelnetwork/pastel/src/pastel-cli -datadir=/home/nd/pastelnetwork/mock-networks/regtest/unzipped/node13 generate 100 

# echo "send 1000 coins to each node"

# for i in $(seq 0 18); do 

# done

