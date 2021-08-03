# GoNode

`gonode` contains two main applications, `walletnode` and `supernode`. These are designed to register NFTs (artworks) on the [Pastel](http://pastel.wiki/en/home) blockchain. Neither `walletnode` nor `supernode`  interact directly with the blockchain; instead, they use Pastel's RPC API to communicate with Pastel's [cNode](https://github.com/pastelnetwork/pastel) which handles the blockchain itself.

### Main Apps

[`walletnode`](walletnode/README.md) is a command line client application (i.e., there is no graphical interface); the main purpose of the app is to provide a REST API on the client side.

[`supernode`](supernode/README.md) is a server application for `walletnode` and actually does all the work of registering and searching for NFTs (artworks).

### Tools

[`pastel-api`](tools/pastel-api/README.md) simulates the Pastel RPC API which is provided by [cNode](https://github.com/pastelnetwork/pastel) and is used for local testing of walletnode and supernode.

### Local testing
1. Clone [`regtest-network`](https://github.com/pastelnetwork/mock-networks)
2. `cd <your-path>/mock-networks/regtest && mkidr unzipped && cd unzipped && tar -xvf ../node.tar.gz`.
3. `cp <your-path>/regtest/{masternode.conf <your-path>/regtest/unzipped/node14/regtest/`.
4. `cp <your-path>/regtest/{start.py, start_masternodes.sh} r <your-path>/regtest/`.
6. `cd <your-path>/regtest && ./start.py`.  
7. `./start_masternodes.sh` and retry if cNode(s) is not ready.
8. Verify if all masternode are `ENABLE` using `<your-path>/pastel-cli -datadir=<your-path>/regtest/unzipped/node14 masternode list`,
9. Generate some block to drive the connections between masternode(s) `<your-path>/pastel-cli -datadir=<your-path>/regtest/unzipped/node13 generate 1`
9. Verify top ranks masternode contains nodes whose `extAddr` is `localhost:{4444:4445:4446}`: `<your-path>/pastel-cli -datadir=<your-path>/regtest/unzipped/node14 masternode top`. Those nodes are the ones supernode(s) will connect to.
10. Verify masternode{0,1,2} also have the same top node lis: `<your-path>/pastel-cli -datadir=<your-path>/regtest/unzipped/node{0,1,2} masternode top`.
11. Start walletnode using: `GRPC_TRACE=tcp,http,api GRPC_VERBOSITY=debug WALLETNODE_DEBUG=1 ./walletnode/walletnode --log-level debug -c ./walletnode/examples/configs/localnet.yml --swagger 2>&1 | tee wallet-node.log`.
12. Start supernode(s) using: `SUPERNODE_DEBUG=1 LOAD_TFMODELS=0 ./supernode/supernode --log-level debug -c ./supernode/examples/configs/localnet-{4444,4445,4446}.yml  2>&1 | tee sn{1,2,3}.log` 
13. Create an artist pastel Id: `<your-path>/pastel-cli -datadir=<your-path>/regtest/unzipped/node14 pastelid newkey passphrase`. The passphrase of artist pastel id is `passphrase`,
14. Create an address: `<your-path>/pastel-cli -datadir=<your-path>/regtest/unzipped/node14 getnewaddress`.
15. Register the artist with network:  `<your-path>/pastel-cli -datadir=<your-path>/regtest/unzipped/node14 tickets register id <artist-pastelId> "passphrase" <artist-address>`.
16. Send coin to artist's adress: `<your-path>/pastel-cli -datadir=<your-path>/regtest/unzipped/node13 sendtoaddress <artist-addr> <amount> "" "" true`.
17. Open your browser at `localhost:8080/swagger`.
18. Upload some image.
19. Register task, change `artist-pastelId`, `spendable-address`, `passphrase` to "passphrase" and `image-id` using the result of upload REST call.
20. Observe the log and use the command to generate coin to drive the flow. The coin need to be generated at following points.
- Walletnode sends `preburnt-txid` to supernode(s).
- Walletnode recieved `reg-art-txid` from supernode(s).
- Walletnode recieved `reg-act-txid` from the network.
- Each step generate 10 block is safe to drive the flow.

**Notes**:
- Walletnode's cNode is node 14.
- Miner is node 13.
- Node{0,1,2} are Node of supernoder{4444,4445,4446} accordingly.

**Troubleshooting**:
- If masternode(s) are `EXPIRE` then do step `7` again.
- If your masternode(s) failed to started due to metadab then remove these folder(s) `~/.pastel/supernode/metadb-444*`
- If your masternode(s) failed to started due to kamedila then remove these folder(s) `~/.pastel/supernode/p2p-localnet-600*`
