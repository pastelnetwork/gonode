package mock

import (
	"fmt"

	"github.com/pastelnetwork/gonode/integration/fakes/common/testconst"
)

var masterNodeExtraResp = `{
	"8edc210250df035a748ef15da1d10c330049c1ce3a76c9987dab4a5f14da21a8-0": {
	  "extAddress": "sn-server-1:14444",
	  "extP2P": "sn-server-1:14445",
	  "extKey": "jXZX9sH6Fp5EUhqKQJteiivbNyGhsjPQKYXqVUVHRDTpH9UzGeEuG92c2S5tMTzUz1CLARVohpzjsWLUyZXSGN",
	  "extCfg": ""
	},
	"0acba0281ee5d753d8fede2bde5e9e1c05351258465365a1f25d79f3808c756a-1": {
	  "extAddress": "sn-server-2:14444",
	  "extP2P": "sn-server-2:14445",
	  "extKey": "jXZFvrCSNQGfrRA8RaXpzm2TEV38hQiyesNZg7joEjBNjXGieM8ZTyhKYxkBrqdtTpYcACRKQjMLQrmKNpfrL8",
	  "extCfg": ""
	},
	"13313ea422bbe97b746782a1ef3770202a5496ee9c239b831062c8aa03a3eed7-1": {
	  "extAddress": "sn-server-3:14444",
	  "extP2P": "sn-server-3:14445",
	  "extKey": "jXXaQm7VuD4p8vG32XyZPq3Z8WBZWdpaRPM92nfZZvcW5pvjVf7RsEnB55qQ214WzuuXr7LoNJsZZKth3ArKAU",
	  "extCfg": ""
	},
	"b2734441a24ecfe8a699e62f3711311100c53cb534f40e6c885948515e410525-1": {
		"extAddress": "sn-server-4:14444",
		"extP2P": "sn-server-4:14445",
		"extKey": "jXYegY26dJ1auzTj5xy1LSUHy3vEn91pMjQrMwAb97nVzH1859xBzbdiU1RsAhJiQ6Z8VbJwQeH5pwGbFAAaHc",
		"extCfg": ""
	  }
  }`
var masterNodesTopResp = `{
	"158978": [
	  {
		"rank": "1",
		"IP:port": "snserver-1:19933",
		"protocol": 170008,
		"outpoint": "8edc210250df035a748ef15da1d10c330049c1ce3a76c9987dab4a5f14da21a8-0",
		"payee": "tPkUm2pJ32LiHaLyugU962LdDSDg87XDZJ2",
		"lastseen": 0,
		"activeseconds": -1640641036,
		"extAddress": "snserver-1:14444",
		"extP2P": "snserver-114445",
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
		"IP:port": "sn-server-4:19933",
		"protocol": 170008,
		"outpoint": "b2734441a24ecfe8a699e62f3711311100c53cb534f40e6c885948515e410525-1",
		"payee": "tPo1skbEb46hq9WWjKyamLhoc1pSmDoa6xo",
		"lastseen": 0,
		"activeseconds": -1640641036,
		"extAddress": "sn-server-4:14444",
		"extP2P": "sn-server-4:14445",
		"extKey": "jXYegY26dJ1auzTj5xy1LSUHy3vEn91pMjQrMwAb97nVzH1859xBzbdiU1RsAhJiQ6Z8VbJwQeH5pwGbFAAaHc",
		"extCfg": ""
	  }
	]
  }`

var networkStorageResp = `{"networkfee": 50}`
var ticketsFindID = `{
	"height": 109374,
	"ticket": {
	  "address": "",
	  "id_type": "masternode",
	  "outpoint": "",
	  "pastelID": "",
	  "pq_key": "",
	  "signature": "040d8393b21c5d8cdbb31cb4e7b4eaa40116ddbad5bfe59488f645c229898a86ce10915b9f5d56f2f98eb7dea50e86387c91710628d0e147808b698510e1b8533f1b733a4f44790ab1ef4ddfc43421eeb5a609d562a41fc446045ea453d7543932007d852342c96856a437654c4317d03d00",
	  "timeStamp": "1634664819",
	  "type": "pastelid",
	  "version": 1
	},
	"txid": "70f80e88289fc0cf2c1fb6a8dbc0090f6c8745df7f799fb506696d7dbccf2670"
  }
  `

var blockVerboseResponse = `{
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
var signResp = `{
	"signature": "DOuI5aeRdNYeKG/eAqIy4SzDSYvTF6biWyEWdwNidJRhKj2N7Ecm6ljIZyXOC2LmpIQw/APw9USADvns1eYLKoxgBB+zkPFtWewSGjasDDGlaw2lk0SVT4HxEVhzABZ5qyzMlVs4VPbaR5NMhs9JAgkA"
  }
  `

var verifyResp = `{
	"verification": "OK"
  }`

var ticketsFindIDArtistResp = `{
	"height": 150940,
	"ticket": {
	  "address": "tPoXb35ABswxqMCw3gD8K7PvzRTcZLdpNEY",
	  "id_type": "personal",
	  "pastelID": "jXY1wJkRFt4hsPn6LnRqUtoRmBx5QTiGcbCXorKq7JuKVy4Zo89PmE8BoGjyujqj6NwfvfGsxhUH2ute6kW2gW",
	  "pq_key": "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",
	  "signature": "a95fd4ac5bff934404b0ae3ad8e04a2ab65993f9ff4004955298d41501a2ae50ce1a6a081858131fac2cd49713cefbef944f1629589cb4b30098b5cb4310b69963b2510da930addd7ab1efaeeaaaa9764a08ee41f4db98fc883f3c9105000c37443116201074855a599939036a7a49bf2300",
	  "timeStamp": "1641492657",
	  "type": "pastelid",
	  "version": 1
	},
	"txid": "896950d860eaf408e76a1a153deff80a7cda9e76291e5085060634e30b145c6a"
  }`

var opStatusResp = `[
	{
	  "id": "opid-0ded54d5-f87a-479d-bcc7-e77866e51409",
	  "status": "success",
	  "creation_time": 1643254712,
	  "result": {
		"txid": "79d01d65154ae5605e2e754c11eda56434c67eb60ad027d66127ab99cdbd0e27"
	  },
	  "execution_secs": 0.008639777,
	  "method": "z_sendmanywithchangetosender",
	  "params": {
		"fromaddress": "tPoXb35ABswxqMCw3gD8K7PvzRTcZLdpNEY",
		"amounts": [
		  {
			"address": "tPpasteLBurnAddressXXXXXXXXXXX3wy7u",
			"amount": 50.0
		  }
		],
		"minconf": 1,
		"fee": 0.1
	  }
	}
  ]`

var storageFeeResp = `{
	  "totalstoragefee":20
  }
  `
var actionFeeResp = `{
	"sensefee":30,
	"cascadefee":40
}
`

var getRawTxResp = `{
	"confirmations":12,
	"txid":""
}
`

var getNFTRegisterResp = `{
	"txid":"c4b1fc370983c7409ec58fcd079136f04efe1e1c363f4cd8f4aff8986a91ef06"
}
`
var getActionRegisterResp = `{
	"txid":"c4b1fc370983c7409ec58fcd079136f04efe1e1c363f4cd8f4aff8986a91ef06"
}
`
var getActRegisterResp = fmt.Sprintf(`{"txid":"%s"}`, testconst.RegActTxID)
