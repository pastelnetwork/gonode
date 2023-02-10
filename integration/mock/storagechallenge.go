package mock

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/integration/helper"
	"github.com/pastelnetwork/gonode/pastel"
)

func getSCRegTickets(oti []byte) (map[string][]byte, []byte) {
	resp := make(map[string][]byte)
	tix := pastel.RegTickets{}

	txid := "b4b1fc370983c7409ec58fcd079136f04efe1e1c363f4cd8f4aff8986a91ef09"
	ticket := pastel.RegTicket{
		TXID: txid,
		RegTicketData: pastel.RegTicketData{
			TotalCopies:   5,
			CreatorHeight: 1,
			StorageFee:    5,
			NFTTicketData: pastel.NFTTicket{
				BlockNum: 1,
				Copies:   2,
				AppTicketData: pastel.AppTicket{
					CreatorName:   "Alan Majchrowicz",
					NFTTitle:      "nft-1",
					NFTSeriesName: "nft-series",
					RQIDs: []string{
						"4aYBeeCEoG7SSnZzXWjq8Ts2oXQ6wW18DGNiEACJJrc4",
					},
					RQOti: oti,
					Thumbnail1Hash: []byte{
						100, 96, 227, 164, 165, 224, 135, 208, 183, 238, 17, 48, 154, 166, 40,
						20, 0, 215, 158, 151, 198, 123, 115, 153, 162, 81, 118, 200, 22, 153, 174, 130,
					},
					Thumbnail2Hash: []byte{249, 23, 138, 33, 126, 119, 199, 251, 112, 210, 244, 236, 144, 50,
						149, 172, 64, 148, 230, 63, 44, 221, 144, 63, 173, 131, 97, 25, 227, 105, 181, 127},
					DataHash: []byte{237, 203, 202, 207, 43, 36, 172, 173, 251, 169, 72, 216, 216, 220, 47,
						235, 33, 171, 187, 188, 199, 189, 250, 43, 39, 154, 14, 144, 51, 135, 12, 68},
					DDAndFingerprintsIDs: []string{ddAndFpFileHash},
				},
			},
		},
	}

	artTicketBytes, _ := pastel.EncodeNFTTicket(&ticket.RegTicketData.NFTTicketData)

	ticket.RegTicketData.NFTTicket = artTicketBytes

	tix = append(tix, ticket)
	t, _ := json.Marshal(ticket)
	resp[txid] = t

	tlist, _ := json.Marshal(tix)

	return resp, tlist
}

var blockVerboseResponseSC = `{
	"hash": "0000034f28e923c3b6393d71c9dea286c67fc33d55d3a0a607b3df40e0d34c95",
	"confirmations": 3,
	"size": 2426,
	"height": 160647,
	"version": 4,
	"merkleroot": "77747bc465eaf2c5aace34f16f0f87d2f729a699f5f8af105c72c3d6af2fb91b",
	"finalsaplingroot": "2820f21343784bbda860087e8d2d50b38947918ccbe108dfbdf0ddde61891da2",
	"tx": [
	  "bdff9c50f269af9c48c9d50d1d155d94042a2d07f6630ea6a05e64b7fc506e58",
	  "a6c57aadc6b5d001bbbf6be98c429169ea2630411ee9d4e8d3e9df01fc873718"
	],
	"time": 1642960828,
	"nonce": "0000000000000000000000000003000000000000642b00c000000000f5ffff17",
	"solution": "009442b6470a44abd7ca486c9cfcd6728b9d7d95280fb7b10dda54bb5fee9a1589e14c5af1d57bdedd9b0473e0fc00cc5ef46639b0bae0c463e8664aac6ee506a203d6495fca7329daf2910dc94d9712ef5fa48d0295ad84c21ec3ff7808a41074fd66254e7eb3e25f0459f14d63e6100151a863c3a6fb3a96cf5bb954ff0ade4a43b9c5816ecdd0d5b9e0d5b785bf71db679d204474493f0967e37e581207e268e8515cf5d7acf7038ba6500b2f1df3d1fc02f519c11eda9f7439b93b0d34edd848082f9878cec1e810de3fe8dea00eccf41a3e47912e20411fbbd9238361a0802ab24d9877f233b13df78e8eb6f0fb36063cd939163fb63dff24ae0879f6b8b8461fe94a2b42a40ba543ab75779f11910c6bb7a53808b9254cafb4a116c9d492be375c40481a74948443db05276aa0a50bbcedd306d9a45f08b34cab32c53cd9b67fcc7f393fb45957925e511f4e2900b67f62008f7863810cf204a6ec6f0ae23c5c0b0c05fd043f2d85558c76ab638a10e9a5910c459bb0b00d5e01ad35442d83dacae0d8c1dbf4921966739cec4e7edd35871f09c3dfb126de38f58991e4add4d68407b0aa4abe50628198ad589a2bd768f2351e183e0e181afb35214a08018261c4a78ad8ef71a72f2e20931da3363bc25181df73f9d3d260ea9279831c74a9d048ca835650a9363d635ed7846651139bcbb7dfc6d200e65e35fa0788a763dd4230e31ed614922cf7de4a289dbb812b2d9ce5bc384672a0497b85f513bdb2ad040d5ffb13be0403fd13518d3c57b86dc4f31616aa20a15d3211cc219ce83ac59fd25989d97c79ef85ac05aa438f595a11773435f2469d7a313f76db1eae8b13b5348e101fdca1a973f47e77d59efdcebdb2e2b64d8bbfe49a9d0315497fe6fb0c5b173b1416dc2edb7696acb80a297a05f87b7bceec7852cbbd60be76d7011ff6787b2eabbf9d3337decd7c0fb2139b944bcf060df11dcda25d4fca87c2442fe5647d770312b7e608c80e3a53891d23bdf752ae8a9b0a19a282f8923f10ab6c19fedc9f995048c20e6a2bdb4a2ef712fcbf08d7077e8652b8d11f7412b1df3eb548bd362f889e489d0dd83d13e12efbf3c86ba74bff32f0325fba931c4989e6d86e9ab77b76a5395641aba57e4dd282c030227b8cf8602c3d43a7f43d59c38535fa2a7965a4012efd42620d3c78b516b0bb75096a40c49127dd87066af9c2957259e5edc452e4b84906b5a93c58241e0edbb67ec8ab6ea996d4b134ce8fb7a2f2c5db74343729ba2d3ee2e4cdf040d5464adf77abb4497e762c0863a16da6d2bd75a0732170906946a1202d9bd6983dfc8fb0daddfd23f46d150e1037e9c9b0aa590c050a6ce9dcf5ce82bb458b57cd8675ce82b655fb61781306b35d8de42343c8c20567e2fb2f4256e4b98f39012925d80c07d814cb61f133232edcf0c9d02eea8b1eae711ae1761fe9e0b2865b903ec0021d58748ab01b77f6829cdf278bc961a1fbf855b221a42fcff9c61d08b61b486133f5fe27141c48c314fdba27f9709d0a62e910cecbc5d6dbf5a8c48cd804f6b9beb72dbc4ada62ee5a16688f6ea4451ef6bb3b5d7850aeaa120b902e998d9f7863895db2a5ba4b269d82c15506d620009bd2cc8c6b691f42d2bb30fddfe374a7fc55430362ca9a41f1d06990431638bddba409f6db3174bb1e0f8cd8e41c021bb19c62c3f5da92d68d461cf8051d3e65e059960e3cea29c2f66e773909ecd133711c24a7f4d1b5566c99253f04632fed53764f91debd3b04e27591d8e40dbb507391c763380e8d1289df0d310643ad011a1d4b9dc2718a00f8e7ad77b51dddd8db07a174c3a54d0f99162c718c9dd747ecb0bf2a5baa1f102788388f4d9d98ed43b58b3a47da55e3dada6c",
	"bits": "1e0ce6f3",
	"difficulty": 40635.64502841329,
	"chainwork": "0000000000000000000000000000000000000000000000000000003efaaacaf1",
	"anchor": "59d2cde5e65c1414c32ba54f0fe4bdb3d67618125286e6a191317917c812c6d7",
	"valuePools": [
	  {
		"id": "sprout",
		"monitored": true,
		"chainValue": 0.00000,
		"chainValuePsl": 0,
		"valueDelta": 0.00000,
		"valueDeltaPsl": 0
	  },
	  {
		"id": "sapling",
		"monitored": true,
		"chainValue": 0.00000,
		"chainValuePsl": 0,
		"valueDelta": 0.00000,
		"valueDeltaPsl": 0
	  }
	],
	"previousblockhash": "000008c298acee1a1adccd7c5071ef8c19a4f573d030e944ed166e91ed75a36f",
	"nextblockhash": "0000047a9f7ed99f34711ecbad76ba76dd6e00883a9752f765e59a5923c5bcd5"
  }
  `
var masterNodeExtraRespSC = `{
	"8edc210250df035a748ef15da1d10c330049c1ce3a76c9987dab4a5f14da21a8-0": {
	  "extAddress": "snserver-1:14444",
	  "extP2P": "snserver-1:14445",
	  "extKey": "jXZX9sH6Fp5EUhqKQJteiivbNyGhsjPQKYXqVUVHRDTpH9UzGeEuG92c2S5tMTzUz1CLARVohpzjsWLUyZXSGN",
	  "extCfg": ""
	},
	"0acba0281ee5d753d8fede2bde5e9e1c05351258465365a1f25d79f3808c756a-1": {
	  "extAddress": "snserver-2:14444",
	  "extP2P": "snserver-2:14445",
	  "extKey": "jXZFvrCSNQGfrRA8RaXpzm2TEV38hQiyesNZg7joEjBNjXGieM8ZTyhKYxkBrqdtTpYcACRKQjMLQrmKNpfrL8",
	  "extCfg": ""
	},
	"13313ea422bbe97b746782a1ef3770202a5496ee9c239b831062c8aa03a3eed7-1": {
	  "extAddress": "snserver-3:14444",
	  "extP2P": "snserver-3:14445",
	  "extKey": "jXXaQm7VuD4p8vG32XyZPq3Z8WBZWdpaRPM92nfZZvcW5pvjVf7RsEnB55qQ214WzuuXr7LoNJsZZKth3ArKAU",
	  "extCfg": ""
	},
	"b2734441a24ecfe8a699e62f3711311100c53cb534f40e6c885948515e410525-1": {
		"extAddress": "snserver-4:14444",
		"extP2P": "snserver-4:14445",
		"extKey": "jABegY26dJ1auzTj5xy1LSUHy3vEn91pMjQrMwAb97nVzH1859xBzbdiU1RsAhJiQ6Z8VbJwQeH5pwGbFAAaHc",
		"extCfg": ""
	  },
	"17761cf79f1ef0be0618759944bae28be7a8f2d41b8fd2e5b7752f76c0018fdf-0": {
		"extAddress": "snserver-5:14444",
		"extP2P": "snserver-5:14445",
		"extKey": "jXa7Ypp5fueoTboLYxAbH2ebaGLiQBRDF2u2PFAMYYGRjSMsEzn5sgSgyAy5xYpEfcMUqaUY6xZ8Qcxnn1pWAf",
		"extCfg": ""
	},
	"7a41c8674e614c5c09080c2f9929e4d87d0c04494b35359f878a47e3482d44f2-1": {
		"extAddress": "snserver-6:14444",
		"extP2P": "snserver-6:14445",
		"extKey": "jXXQD8CjECcKknAUf5NHcCua8aqHhbq4vP5EyE6kCDvD6EWpnbjaUe8VH9ZVpvj8n85hkMtDZyPqrwJXPwYUBi",
		"extCfg": ""
	},
	"8da821743ceb5932d7a10facee121dbfed6e708e76ab1d2dd15c315577be70c0-1": {
		"extAddress": "snserver-7:14444",
		"extP2P": "snserver-7:14445",
		"extKey": "jXYqgELnZQUSC2D5ZKPzVxaJwDLKC2NAzTzvEkR5EXU3od7xDJH9JgsMWvU6FoYQaBVP489A1y7qw756m8Er1S",
		"extCfg": ""
	}
  }`
var masterNodesTopRespSC = `{
	"160647": [
	  {
		"rank": "1",
		"IP:port": "snserver-1:19933",
		"protocol": 170008,
		"outpoint": "8edc210250df035a748ef15da1d10c330049c1ce3a76c9987dab4a5f14da21a8-0",
		"payee": "tPkUm2pJ32LiHaLyugU962LdDSDg87XDZJ2",
		"lastseen": 0,
		"activeseconds": -1640641036,
		"extAddress": "snserver-1:14444",
		"extP2P": "snserver-14445",
		"extKey": "jXZX9sH6Fp5EUhqKQJteiivbNyGhsjPQKYXqVUVHRDTpH9UzGeEuG92c2S5tMTzUz1CLARVohpzjsWLUyZXSGN",
		"extCfg": ""
	  },
	  {
		"rank": "2",
		"IP:port": "snserver-2:19933",
		"protocol": 170008,
		"outpoint": "0acba0281ee5d753d8fede2bde5e9e1c05351258465365a1f25d79f3808c756a-1",
		"payee": "tPXkp3hxWFsWJ9PM33sfj53MqXEtKmALm1t",
		"lastseen": 0,
		"activeseconds": -1640641036,
		"extAddress": "snserver-2:14444",
		"extP2P": "snserver-2:14445",
		"extKey": "jXZFvrCSNQGfrRA8RaXpzm2TEV38hQiyesNZg7joEjBNjXGieM8ZTyhKYxkBrqdtTpYcACRKQjMLQrmKNpfrL8",
		"extCfg": ""
	  },
	  {
		"rank": "3",
		"IP:port": "snserver-3:19933",
		"protocol": 170008,
		"outpoint": "13313ea422bbe97b746782a1ef3770202a5496ee9c239b831062c8aa03a3eed7-1",
		"payee": "tPbhJWPk2YkcVizX7yGr1AcjBVQwt5qXNji",
		"lastseen": 0,
		"activeseconds": -1640641036,
		"extAddress": "snserver-3:14444",
		"extP2P": "snserver-3:14445",
		"extKey": "jXXaQm7VuD4p8vG32XyZPq3Z8WBZWdpaRPM92nfZZvcW5pvjVf7RsEnB55qQ214WzuuXr7LoNJsZZKth3ArKAU",
		"extCfg": ""
	  },
		{
		"rank": "4",
		"IP:port": "snserver-4:19933",
		"protocol": 170008,
		"outpoint": "b2734441a24ecfe8a699e62f3711311100c53cb534f40e6c885948515e410525-1",
		"payee": "tPo1skbEb46hq9WWjKyamLhoc1pSmDoa6xo",
		"lastseen": 0,
		"activeseconds": -1640641036,
		"extAddress": "snserver-4:14444",
		"extP2P": "snserver-4:14445",
		"extKey": "jABegY26dJ1auzTj5xy1LSUHy3vEn91pMjQrMwAb97nVzH1859xBzbdiU1RsAhJiQ6Z8VbJwQeH5pwGbFAAaHc",
		"extCfg": ""
	  },
	  {
		"rank": "5",
		"IP:port": "snserver-5:19933",
		"protocol": 170008,
		"outpoint": "17761cf79f1ef0be0618759944bae28be7a8f2d41b8fd2e5b7752f76c0018fdf-0",
		"payee": "tPeFBz1WC5VY6mtPV5Tn1vTtQKfeMnvsS85",
		"lastseen": 0,
		"activeseconds": -1640641036,
		"extAddress": "snserver-5:14444",
		"extP2P": "snserver-5:14445",
		"extKey": "jXa7Ypp5fueoTboLYxAbH2ebaGLiQBRDF2u2PFAMYYGRjSMsEzn5sgSgyAy5xYpEfcMUqaUY6xZ8Qcxnn1pWAf",
		"extCfg": ""
	  },
	  {
		"rank": "6",
		"IP:port": "snserver-6:19933",
		"protocol": 170008,
		"outpoint": "7a41c8674e614c5c09080c2f9929e4d87d0c04494b35359f878a47e3482d44f2-1",
		"payee": "tPSeVpnQAgoH3Z5Wzz71pBRdGESCvvR1EbX",
		"lastseen": 0,
		"activeseconds": -1640641036,
		"extAddress": "snserver-6:14444",
		"extP2P": "snserver-6:14445",
		"extKey": "jXXQD8CjECcKknAUf5NHcCua8aqHhbq4vP5EyE6kCDvD6EWpnbjaUe8VH9ZVpvj8n85hkMtDZyPqrwJXPwYUBi",
		"extCfg": ""
	  },
	  {
		"rank": "7",
		"IP:port": "snserver-7:19933",
		"protocol": 170008,
		"outpoint": "8da821743ceb5932d7a10facee121dbfed6e708e76ab1d2dd15c315577be70c0-1",
		"payee": "tPgJre49g2heWpxvD3tY1RUhbWqBgpCbVk8",
		"lastseen": 0,
		"activeseconds": -1673572702,
		"extAddress": "snserver-7:14444",
		"extP2P": "snserver-7:14445",
		"extKey": "jXYqgELnZQUSC2D5ZKPzVxaJwDLKC2NAzTzvEkR5EXU3od7xDJH9JgsMWvU6FoYQaBVP489A1y7qw756m8Er1S",
		"extCfg": ""
	  }
	]
  }`

func GetRqIDFile() *helper.RawSymbolIDFile {
	rqfile := &helper.RawSymbolIDFile{
		ID:        "7aa4746a-450d-4cbb-a5a1-9bef7b0eb58d",
		PastelID:  "jXZvhdVoQ2q2WfaqL2nRCVNZn5hKgGDYGpkaQsjco9AgQA5MH1he4QktDN5RP483qyN17SPFs34o73tjLhfnxs",
		BlockHash: "0000034f28e923c3b6393d71c9dea286c67fc33d55d3a0a607b3df40e0d34c95",
	}

	rqfile.SymbolIdentifiers = []string{`38uydDwESDz37Xg7Hhn7z2yhCTfjmhLcAHgdpJoWMCHP`, `3vgZe9zyP3caEwKJcLZ78wHWQfByi6aW5E6EbGmYWvM1`,
		`893qE6a4QtrYDFzgJVivDq39W2bcGKfeJcnWQjxtoJzC`, `Enmj4KkuBEGeTx4puPdhQTU3gYxL4Rf7PFnmtFhQnpmN`, `6AU9Ti8VfQ8ojTRN4ctNACJ7C8cJksvYruko9pC1rh24`,
		`AqGutQGxyUw1nD4E4XsKaLnaeZpcQ8P89dzA1PZprt5r`, `DSFV5iFrfDz9ciMudoVYbry35EoUt4LLqP9GBMdpGRJe`, `FEzCE2kcNXxPezgpkftk7n2p6tDgbpSVso5pgtBQz5Zk`,
		`CUwccZ6UKeA1CbF8YeyqN8cX9dEn8dEcj7Kw5dcyEnnB`, `MoqRvvKYL3NoNAXtS8e1X1FwLSit76YPjsnkdQZAfK3`, `7PaDagNq6Dm1UBhC9YAb7hj74cpTbjQtu4GnZBhK3M3F`,
		`5vemfsETYz7xEgTDAgyq6wXmxkNeehw8TP6qrMQVP1WZ`, `44ECnUJHJzvs7MM9LzXLqRtXHPSNL63mS8zQ3xbV7MjR`, `6N25xCJM9z1BMYdSPP7vf293wt1SNB7XaifNbNNd1ang`,
		`2dWHDqpbSTtwDg9KjKdZxTA2hqHgut6tCktRPcrzep5V`, `CaHU26Z3XUDMRUfxf7qnsUJPT7NeobVWSczjASG5oYua`, `Ad1M7xuwsZ3P4JYP5Vxuzn2gcKGC8hXeKtpyoL5MDAf1`,
		`6hdN4PPsSoPRi79Wr4j3utxuWZEAR9vRt8YoqUZZ8E4g`}

	return rqfile
}
