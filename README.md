# GoNode

[![PastelNetwork](https://circleci.com/gh/pastelnetwork/gonode.svg?style=shield)](https://app.circleci.com/pipelines/github/pastelnetwork/gonode)

GoNode comprises two primary applications, `walletnode` and `supernode`. They are engineered to facilitate NFT registration on the [Pastel](https://docs.pastel.network/introduction/pastel-overview) blockchain. Neither `walletnode` nor `supernode` directly interface with the blockchain. Instead, they employ Pastel's RPC API to interact with Pastel's [cNode](https://github.com/pastelnetwork/pastel), which manages the blockchain itself.

## Build

- Begin with [`golang`](https://go.dev/doc/) and install the latest version.
- To generate mocks, install [`mockery`](https://github.com/vektra/mockery).
- For proto generation, install the protoc compiler, protoc-gen-go, and protoc-gen-go-grpc:
```
 go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
 go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
```

Run `make gen-mocks` if want to build and run unit tests

```
make build
```

## Testing

There are multiple methods and tools available to test the functionality of gonode. However, since the entire system relies on cNode and a collection of remote SuperNodes, setting up such an environment can be challenging.

### Unit Tests
```shell
make sn-unit-tests
make wn-unit-tests
```

### Integration Tests

Integration tests ensure API contracts are upheld and the APIs function as anticipated. These tests launch supernode containers and emulate CNode, RaptorQ & Duplication Detection services. These mock services are located in `/fakes`. 

```shell
make integration-tests
```

For more information about Integration Tests, refer to the README in the `/integration` directory.

### Testing on Testnet

1. Clone [Pastelup](https://github.com/pastelnetwork/pastelup).
2. Install walletnode by following the steps provided in the [pastelup README](https://github.com/pastelnetwork/pastelup/blob/master/README.md).
3. Ensure your cNode is fully synced and returns top nodes when you run `./pastel-cli masternodes top`.
4. Create an artist pastel Id: `./pastel-cli pastelid newkey passphrase`. Here, the passphrase for artist pastel id is `passphrase`.
5. Generate an address: `./pastel-cli getnewaddress`. You'll use the generated address in the next step.
6. From another node on the testnet, send some coins to the newly created artist's address: `./pastel-cli sendtoaddress <artist-addr> <amount> "" "" true`.
7. Register the artist with the network: `./pastel-cli tickets register id <artist-pastelId> "passphrase" <artist-address>`. After this step, you can use your `pastel-id` to send requests to the walletnode.

### Testing on RegTest

The [`pastel-api`](tools/pastel-api/README.md) simulates the Pastel RPC API, provided by [cNode](https://github.com/pastelnetwork/pastel), and is used for local testing of the walletnode and supernode.

Follow these steps:

1. Clone [`regtest-network`](https://github.com/pastelnetwork/mock-networks).
2. Unzip node files: `cd <your-path>/mock-networks/regtest && mkdir unzipped && cd unzipped && tar -xvf ../node.tar.gz`.
3. Copy `masternode.conf` to node directory: `cp <your-path>/regtest/masternode.conf <your-path>/regtest/unzipped/node14/regtest/`.
4. Copy `start.py` and `start_masternodes.sh` to regtest directory: `cp <your-path>/regtest/{start.py,start_masternodes.sh} <your-path>/regtest/`.
5. Start python script: `cd <your-path>/regtest && ./start.py`.
6. Execute `./start_masternodes.sh` and retry if cNode(s) are not ready.
7. Verify if all masternodes are `ENABLE` using `<your-path>/pastel-cli -datadir=<your-path>/regtest/unzipped/node14 masternode list`.
8. Generate a block to establish connections between masternodes: `<your-path>/pastel-cli -datadir=<your-path>/regtest/unzipped/node13 generate 1`.
9. Verify top-ranked masternodes contain nodes whose `extAddr` is `localhost:{4444:4445:4446}` using `<your-path>/pastel-cli -datadir=<your-path>/regtest/unzipped/node14 masternode top`. These nodes are the ones supernode(s) will connect to.
10. Verify masternode{0,1,2} also have the same top node list: `<your-path>/pastel-cli -datadir=<your-path>/regtest/unzipped/node{0,1,2} masternode top`.
11. Start walletnode: `GRPC_TRACE=tcp,http,api GRPC_VERBOSITY=debug WALLETNODE_DEBUG=1 ./walletnode/walletnode --log-level debug -c ./walletnode/examples/configs/localnet.yml --swagger 2>&1 | tee wallet-node.log`.
12. Start supernode(s): `SUPERNODE_DEBUG=1 LOAD_TFMODELS=0 ./supernode/supernode --log-level debug -c ./supernode/examples/configs/localnet-{4444,4445,4446}.yml  2>&1 | tee sn{1,2,3}.log`.
13. Create an artist Pastel ID: `<your-path>/pastel-cli -datadir=<your-path>/regtest/unzipped/node14 pastelid newkey passphrase`. The passphrase of the artist Pastel ID is `passphrase`.
14. Create an address: `<your-path>/pastel-cli -datadir=<your-path>/regtest/unzipped/node14 getnewaddress`.
15. Register the artist with the network: `<your-path>/pastel-cli -datadir=<your-path>/regtest/unzipped/node14 tickets register id <artist-pastelId> "passphrase" <artist-address>`.
16. Send coins to the artist's address using this command: `<your-path>/pastel-cli -datadir=<your-path>/regtest/unzipped/node13 sendtoaddress <artist-addr> <amount> "" "" true`.
17. Open a web browser and navigate to `localhost:8080/swagger`.
18. Upload an image.
19. Register a task by changing the `artist-pastelId`, `spendable-address`, and `passphrase` to "passphrase". Use the `image-id` obtained from the result of the image upload REST call.
20. Monitor the logs and use the appropriate command to generate coins as needed to maintain the flow. Coins need to be generated at the following stages:
    - When the Walletnode sends `preburnt-txid` to a Supernode.
    - When the Walletnode receives `reg-art-txid` from a Supernode.
    - When the Walletnode receives `reg-act-txid` from the network.
    - Generating 10 blocks at each step is a safe way to keep the flow going.

**Notes**:
- Walletnode's cNode is Node 14.
- The Miner is Node 13.
- Node{0,1,2} correspond to Supernode{4444,4445,4446} respectively.

**Troubleshooting**:
- If the Masternode(s) status is `EXPIRE`, repeat step `7`.
- If any Masternode(s) failed to start due to metadab, remove the respective folder(s) at `~/.pastel/supernode/metadb-444*`.
- If any Masternode(s) failed to start due to kamedila, remove the respective folder(s) at `~/.pastel/supernode/p2p-localnet-600*`.

## Components

### Walletnode
[`walletnode`](walletnode/README.md) operates as an API server, providing REST APIs on the client-side. The APIs enable users to search, view, and download NFTs, as well as register their own NFTs. Additionally, it allows end-users and other applications to utilize the Pastel Network's duplication detection system for images and a storage system. By default, it operates on port :8080. Configurations can be modified via `walletnode.yml` in the Pastel configuration directory.

It offers the following Workflow APIs. The list below provides a high-level overview. For detailed API information, please refer to the swagger docs at `localhost:8080/swagger`

#### NFT Register & Search
Use this workflow to register a new NFT as an artist/creator. You'll need your PastelID registered on the network. The steps to achieve this are detailed in the Testing Section of this document.

  - `http://localhost:8080/nfts/register/upload`: Upload the image and receive a file_id.
  - `http://localhost:8080/nfts/register`: Initiate NFT register task. It returns a task_id in response.
  - `ws://127.0.0.1:8080/nfts/register/<<task_id>>/state`: Connect to this Websocket URL to monitor the status of the NFT register task. A successful registration will display the Task Completed status. NFT registration can take between 20 and 60 minutes.
  - `http://localhost:8080/nfts/<<task_id>>/history`: View task status history and details.
  - `http://localhost:8080/nfts/download?pid=<<pastel_id>>L&txid=<<tx_id>>`: Download the NFT. Only the current owner can download the original NFT.
  - `ws://localhost:8080/nfts/search?query=<<query>>&creator_name=true&art_title=true`: Search for an NFT. This WS provides various filters. Please consult the swagger docs for details.

#### Cascade Register
Cascade facilitates secure and reliable file storage for other applications or users through Pastel's storage layer.

1. Use `http://127.0.0.1:8080/openapi/cascade/upload` to upload a file and retrieve `file_id` and `estimated_fee`.
2. Initiate the Cascade register task by burning 20% of the `estimated_fee`. Find the burn address in `walletnode.yml` and send 20% of the estimated fee to this address from your address (the one used to register your PastelID).
3. Use `http://localhost:8080/openapi/cascade/start/<<file_id>>` to start the Cascade register task, which will return a `task_id` in the response.
4. Monitor the status of the Cascade register task by connecting to this Websocket URL: `ws://127.0.0.1:8080/openapi/cascade/start/<<task_id>>/state`. A successful register will display a `Task Completed` status. Note that the registration process may take between 20 - 60 minutes.
5. Use `http://localhost:8080/openapi/cascade/<<task_id>>/history` to view the task status history and details.

#### Sense Register
Sense provides other applications or users access to Pastel's state-of-the-art Duplication Detection service in exchange for a small fee.

1. Use `http://127.0.0.1:8080/openapi/sense/upload` to upload a file and retrieve `file_id` and `estimated_fee`.
2. Initiate the Sense register task by burning 20% of the `estimated_fee`. Find the burn address in `walletnode.yml` and send 20% of the estimated fee to this address from your address (the one used to register your PastelID).
3. Use `http://localhost:8080/openapi/sense/start/<<file_id>>` to start the Sense register task, which will return a `task_id` in the response.
4. Monitor the status of the Sense register task by connecting to this Websocket URL: `ws://127.0.0.1:8080/openapi/sense/start/<<task_id>>/state`. A successful register will display a `Task Completed` status. Note that the registration process may take between 20 - 60 minutes.
5. Use `http://localhost:8080/openapi/sense/<<task_id>>/history` to view the task status history and details.

#### Cascade & Sense Download 
- To download the original file, use `http://localhost:8080/openapi/cascade/download?pid=<<pastel_id>>L&txid=<<tx_id>>`. Only the owner can download the file.
- To download the duplication detection results for the provided image, use `http://localhost:8080/openapi/sense/download?pid=<<pastel_id>>L&txid=<<tx_id>>`.

### Supernode
[`supernode`](supernode/README.md) is a server application for `walletnode`. It handles communication with the blockchain, duplication detection server, storage, and other services.

### Hermes
Hermes has two primary roles:

- Cleans up inactive tickets: If Registration or Action (Sense\Cascade) tickets aren't activated after a specified block height (configurable), Hermes will remove the data from the local SQLite Database. This process helps maintain a manageable database size and eliminates unnecessary data.

- Stores fingerprints of Registration and Action ticket files in the fingerprints Database.


## More information

[Pastel Network Docs](https://docs.pastel.network/introduction/pastel-overview)


## Makefile tool
- `make sn-unit-tests`: Execute sn unit tests, including relevant packages.
- `make wn-unit-tests`: Execute wn unit tests, including relevant packages.
- `make integration-tests`: Execute integration tests.
- `make build`: Verify if WN & SN can build successfully.
- `make gen-mocks`: Remove existing mocks and generate new ones for interfaces throughout gonode (mocks are used in unit tests).
- `make gen-proto`: Generate all protobuf models & services.
- `make clean-proto`: Remove all protobuf models & services.

