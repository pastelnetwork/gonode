# GoNode

`gonode` contains two main applications `walletnode` and `supernode` are designed to register NFTs (artworks) in the [Pastel](http://pastel.wiki/en/home) blockchain. Neither `walletnode` nor `supernode`  interact directly with the blockchain, they use pastel RPC API to communcation with [cNode](https://github.com/pastelnetwork/pastel) that handles the blockchain itself.

### Main Apps

[`walletnode`](supdernode/README.md) is a client application, without graphical interface, the main purpose the app is to provide REST API, on the client side.
[`supernode`](walletnode/README.md) figuratively speaking, is a server application for `walletnode`, the app does all the work of registering and searching artworks.


### Tools
[`pastel-api`](tools/pastel-api/README.md) imitates Pastel RPC API which is provided by [cNode](https://github.com/pastelnetwork/pastel). Used for local testing walletnode and supernode.
