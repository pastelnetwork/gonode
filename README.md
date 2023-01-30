# GoNode

[![PastelNetwork](https://circleci.com/gh/pastelnetwork/gonode.svg?style=shield)](https://app.circleci.com/pipelines/github/pastelnetwork/gonode)

`gonode` contains two main applications, `walletnode` and `supernode`. These are designed to register NFTs on the [Pastel](https://docs.pastel.network/introduction/pastel-overview) blockchain. Neither `walletnode` nor `supernode`  interact directly with the blockchain; instead, they use Pastel's RPC API to communicate with Pastel's [cNode](https://github.com/pastelnetwork/pastel) which handles the blockchain itself.

## Getting started
- Get started with [`golang`](https://go.dev/doc/) & install the latest version.
- In order to generate mocks, [`mockery`](https://github.com/vektra/mockery) needs to be installed.
- In order to generate proto, protoc compiler needs to be installed as well as protoc-gen-go & protoc-gen-go-grpc

```
 go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
 go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

```

## Walletnode

[`walletnode`](walletnode/README.md) is an API server that provides REST APIs on the client side. The APIs allow the user to search, view & download nfts as well as enable the users to register their own nfts. In addition, it also allows end-users or other applications to use Pastel Network's duplication detection system for the images as well as a storage system. By default, it runs of port `:8080`. Configurations can be changed through `walletnode.yml` located in pastel configuration directory.

It provides the following Workflow APIs. The list below provides a high-level overview. Please consult swagger docs to see the details of APIs  at `localhost:8080/swagger`

#### NFT Register & Search
 Use this workflow to register a new NFT as an artist\creator. For this, you'll need your PastelID registered in the network. The steps to do that are listed below in this document in *Testing* Section.

  - `http://localhost:8080/nfts/register/upload` to upload the image and get `file_id`
  - `http://localhost:8080/nfts/register` to start NFT register task. It returns `task_id` in response
  - `ws://127.0.0.1:8080/nfts/register/<<task_id>>/state` connect with this Websocket URL to monitor the status of NFT register task. For a successful register, you should see `Task Completed` status. It could take anywhere b/w 20 - 60 minutes for an NFT to register successfully.
  - `http://localhost:8080/nfts/<<task_id>>/history` to see task status history and details. 
  - `http://localhost:8080/nfts/download?pid=<<pastel_id>>L&txid=<<tx_id>>` to download the NFT. Only the current owner can download the original NFT.
  - `ws://localhost:8080/nfts/search?query=<<query>>&creator_name=true&art_title=true` to search for NFT. This WS provides many filters. Please consult swagger docs for details. 

#### Cascade Register
Cascade allows other applications or users to use Pastel's storage layer to store files securely and reliably. 

  1. `http://127.0.0.1:8080/openapi/cascade/upload` to upload the file and get `file_id` and `estimated_fee`. <br/>
  2. You'll need to burn 20% of the `estimated_fee` to start cascade register task. Please get burn address from `walletnode.yml` and send 20% of the estimated fee to that address from your address.(The one you used to register your PastelID). <br/>
  3. `http://localhost:8080/openapi/cascade/start/<<file_id>>` to start Cascade register task. It returns `task_id` in response. <br/>
  4. `ws://127.0.0.1:8080/openapi/cascade/start/<<task_id>>/state` connect with this Websocket URL to monitor the status of Cascade register task. For a successful register, you should see `Task Completed` status. It could take anywhere b/w 20 - 60 minutes for a File to register successfully. <br/>
  5. `http://localhost:8080/openapi/cascade/<<task_id>>/history` to see task status history and details. <br/>
  
#### Sense Register
Sense allows other applications or users to use Pastel's state-of-the-art Duplication Detection service by paying a small fee. 

  1. `http://127.0.0.1:8080/openapi/sense/upload` to upload the file and get `file_id` and `estimated_fee`. <br/>
  2. You'll need to burn 20% of the `estimated_fee` to start sense register task. Please get burn address from `walletnode.yml` and send 20% of the estimated fee to that address from your address.(The one you used to register your PastelID). <br/>
  3. `http://localhost:8080/openapi/sense/start/<<file_id>>` to start Sense register task. It returns `task_id` in response. <br/> 
  4. `ws://127.0.0.1:8080/openapi/sense/start/<<task_id>>/state` connect with this Websocket URL to monitor the status of Sense register task. For a successful register, you should see `Task Completed` status. It could take anywhere b/w 20 - 60 minutes for a File to register successfully. <br/>
  5. `http://localhost:8080/openapi/sense/<<task_id>>/history` to see task status history and details. <br/>

#### Casace & Sense Download 
 - use `http://localhost:8080/openapi/cascade/download?pid=<<pastel_id>>L&txid=<<tx_id>>` to download the original file. owner can download the File.
 - use `http://localhost:8080/openapi/sense/download?pid=<<pastel_id>>L&txid=<<tx_id>>` to download the duplication detection results for the provided image.

## Supernode

[`supernode`](supernode/README.md) is a server application for `walletnode` and its responsible for communication with the blockchain, duplication detection server, storage and other services.

## Bridge

Bridge service provides a "bridge" to walletnode to download files through supernode quickly. In order to connect with supernode, walletnode first needs to do a hand-shake which takes a bit of time. Obviously, this hand-shake, when done after a search or fetch request has been made, would delay the response of Wallenode to the end-user. This is where bridge comes in. It maintains connection with top 10 nodes and as soon as Walletnode requests a file, it instantly downloads them from Supernodes and send it back.

## Hermes

Hermes has two major responsibilities 

- Cleanup the inactive tickets: If Registration or Action (Sense\Cascade) tickets are not activated after a certain block height (configurable) then hermes is supposed to cleanup the data from local SQLite Database. Just as it sounds, This helps keeping the database size small and free from useless data. 

- Stores fingerprints of Registration and Action ticket files in the fingerprints Database

## Integration tests

We use Integration tests to make sure API contracts do not break and APIs work as expected. Integration tests spin up supernode containers and fake CNode, RaptorQ & Duplication Detection services. These fake services live in `/fakes`. For more details on Integration Tests, please see README in `/integration` directory.

## Testing
There are various ways & tools to test the working of gonode. Because the entire system is dependent on cNode and collection of remote SuperNodes. It can be challenging to spin up such and enviornment 
### Testing on testnet
1. Clone [Pastelup](https://github.com/pastelnetwork/pastelup) <br/>
2. Install walletnode by following the steps mentioned in [pastelup README](https://github.com/pastelnetwork/pastelup/blob/master/README.md) <br/>
3. Make sure that your cNode is all synced up and return top nodes when you run `./pastel-cli masternodes top` <br/>
4. Create an artist pastel Id: `./pastel-cli pastelid newkey passphrase`. The passphrase of artist pastel id is `passphrase` in this case. <br/>
5. Create an address: `./pastel-cli getnewaddress`. The generated address shall be used in the next step. <br/>
6. From some other node on testnet, send some coins to the generated artist's adress: `./pastel-cli sendtoaddress <artist-addr> <amount> "" "" true`. <br/>
7. Register the artist with network:  `./pastel-cli tickets register id <artist-pastelId> "passphrase" <artist-address>`. After this, you'll be able to use your `pastel-id` to send requests onto the walletnode. 

### Testing on regtest

[`pastel-api`](tools/pastel-api/README.md) simulates the Pastel RPC API which is provided by [cNode](https://github.com/pastelnetwork/pastel) and is used for local testing of walletnode and supernode.

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
- Walletnode received `reg-art-txid` from supernode(s).
- Walletnode received `reg-act-txid` from the network.
- Each step generate 10 block is safe to drive the flow.

**Notes**:
- Walletnode's cNode is node 14.
- Miner is node 13.
- Node{0,1,2} are Node of supernoder{4444,4445,4446} accordingly.

**Troubleshooting**:
- If masternode(s) are `EXPIRE` then do step `7` again.
- If your masternode(s) failed to started due to metadab then remove these folder(s) `~/.pastel/supernode/metadb-444*`
- If your masternode(s) failed to started due to kamedila then remove these folder(s) `~/.pastel/supernode/p2p-localnet-600*`

## Makefile tool

- `make sn-unit-tests` execute sn unit tests including packages that it uses
- `make wn-unit-tests` execute wn-unit tests including packages that it uses
- `make integration-tests` execute integration tests
- `make build` verify if WN & SN are able to build
- `make gen-mocks` remove existing mocks & generates all mocks for interfaces throughout gonode (mocks are used in unit tests)
- `make gen-proto` generate all protobuf models & services
- `make clean-proto` removes all protobuf models & services

## Contribution 

- Create a new branch from latest master. Branch name should have the designated ticket number. For example, *[PSL-XX]_Fix-xyz*
- Run the following commands to make sure everything is ok. 

```
make build
make sn-unit-tests
make wn-unit-tests
make integration-tests
```
Feel free to use `make gen-mocks` and\or `make gen-proto` if needed. 

- Open up the PR and keep an eye on CircleCI pipleines to make sure everything works. 
- Wait for atleast two approvals before merging. 