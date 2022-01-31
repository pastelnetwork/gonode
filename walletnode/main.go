// Walletnode provides blockchain "client" functionality, with users able to register nft's,
// download nft's, access the Sense API interface, and interact with chain user data.
//
// Procedural information for blockchain interaction, as well as other developer documentation is located
// at https://pastel.wiki/
//
//                                    ┌──────────────┐
//                                    │              │
//                                    │    Main      │
//                                    │              │
//                                    └──────┬───────┘
//                                           │
//                                    ┌──────▼───────┐
//                                    │              │
//                                    │ App.Run()    │
//                                    │  runApp()    │
//     ...................            └──────┬───────┘
//   ...                 .........           │   runApp Creates the Server
//  ..User activity will interact .   ┌──────▼───────┐
//  .  with this API Server      ..   │              │
// ...............................────►   Server     │
//                                    │              │
//                  ┌──────────┬──────┴─────┬───────┬┴─────────────────┐
//                  │          │            │       │                  │
//                  │          │            │       │                  │  Server creates the
//                  │          │            │       │                  │   services layer
//     ┌────────────▼┐ ┌───────▼─┐  ┌───────▼───┐  ┌▼──────────────┐  ┌▼────────────┐
//     │ NFTRegister │ │NFTSearch│  │NFTDownload│  │UserDataProcess│  │SenseRegister│
//     └─────┬───────┘ └─────────┘  └───────────┘  └───────────────┘  └─────────────┘
//           │                 Each service creates and starts its own
//           │           │   Worker pool, which also starts custom Tasks. │
//           │           │     This is the Task layer, as an example      │
//           │           │      we have NftRegisterTask below             │
//           │
//           │                 ┌─────────────────────────────┐
//           │                 │ NftRegisterTask             │
//           └─────────────────►                             │
//                             │ Controls all interaction wit│   Tasks manage procedural blockchain logic, but
//                             │ the mesh network of nodes   │   most of the actual communication
//                             │ necessary to register an    │   logic happens at a lower layer,
//                             │ NFT on the blockchain.      │   with a connection Session calling
//                             └─┬───────────────────────────┘   a gRPC function, e.g. SendRegMetadata
//                               │
//          ┌────────────────────▼──┐                 ┌────────────────────────────┐
//          │  Session Functions    ├────────────────►│ gRPC Functions             │
//          │  e.g. SendRegMetadata │                 │ e.g. SendRegMetadata       │
//          └───────────────────────┘                 └──────┬──────────────────▲──┘
//                                                           │                  │
//                                                   ********▼***********************
//                                                   ** Supernode Mesh Network      ***
//                                                    *                               **
//                                                   ** Provides  blockchain "Server" **
//                                                   *  functionality by handling     **
//                                                   ** transactions.             ******
//                                                    *****************************

package main

import (
	"os"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/walletnode/cmd"
)

const (
	debugModeEnvName = "WALLETNODE_DEBUG"
)

var (
	debugMode = sys.GetBoolEnv(debugModeEnvName, false)
)

func main() {
	defer errors.Recover(log.FatalAndExit)

	//configuration in here, app.Run will actually run NewApp's runApp function.
	app := cmd.NewApp()
	err := app.Run(os.Args)

	log.FatalAndExit(err)
}

func init() {
	log.SetDebugMode(debugMode)
}
