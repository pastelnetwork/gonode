# GoNode

`gonode` contains two main applications, `walletnode` and `supernode`. These are designed to register NFTs (artworks) on the [Pastel](http://pastel.wiki/en/home) blockchain. Neither `walletnode` nor `supernode`  interact directly with the blockchain; instead, they use Pastel's RPC API to communicate with Pastel's [cNode](https://github.com/pastelnetwork/pastel) which handles the blockchain itself.

### Main Apps

[`walletnode`](supdernode/README.md) is a command line client application (i.e., there is no graphical interface); the main purpose of the app is to provide a REST API on the client side.
[`supernode`](walletnode/README.md) is a server application for `walletnode` and actually does all the work of registering and searching for NFTs (artworks).

### Tools

[`pastel-api`](tools/pastel-api/README.md) simulates the Pastel RPC API which is provided by [cNode](https://github.com/pastelnetwork/pastel) and is used for local testing of walletnode and supernode.
