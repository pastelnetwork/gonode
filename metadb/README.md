# metadb

`metadb` is Metadata Layer, the one that store these following data in RQLite cluster database:
- User specified data (userdata) --> Ready
- Supernode specified data. --> Under Development
- It also have responsible to calculate walletnode and supernode reputation scores, and take action for the low reputation user/walletnode/supernode.--> Under Development

## Quick Start

 
1. Build supernode, walletnode, tools/pastel-api with command:

>cd "target  dir"

>go build ./


2. The application running path depends on the OS:

* MacOS `~/Library/Application Support/Pastel`

* Linux `~/.pastel`

* Windows (>= Vista) `C:\Users\Username\AppData\Roaming\Pastel`

* Windows (< Vista) `C:\Documents and Settings\Username\Application Data\Pastel`

  
3. Testing in local

- Start 3 Supernodes:

>SUPERNODE_DEBUG=1 ./supernode --log-level debug -c ./examples/configs/localnet-4444.yml

>SUPERNODE_DEBUG=1 ./supernode --log-level debug -c ./examples/configs/localnet-4445.yml

>SUPERNODE_DEBUG=1 ./supernode --log-level debug -c ./examples/configs/localnet-4446.yml

  

- Start Walletnode:

>WALLETNODE_DEBUG=1 ./walletnode --log-level debug -c ./examples/configs/localnet.yml --swagger

  
- Start pastel api tool:

>PASTEL_API_DEBUG=1 ./pastel-api --log-level debug

  

- Running swagger api:

>http://localhost:8080/swagger/#/userdatas

  

Execute the swagger api call using the info provided by swagger api, there is following api:

- http://localhost:8080/userdatas/create (POST)

- http://localhost:8080/userdatas/update (POST)

- http://localhost:8080/userdatas/<pastelid> (GET)

- http://localhost:8080/swagger/#/userdatas/userdatas%23setUserFollowRelation (POST)

- http://localhost:8080/swagger/#/userdatas/userdatas%23getFollowees (GET)

- http://localhost:8080/swagger/#/userdatas/userdatas%23getFollowers (GET)

- http://localhost:8080/swagger/#/userdatas/userdatas%23getFriends (GET)

- http://localhost:8080/swagger/#/userdatas/userdatas%23getUsersLikeArt (GET)

- http://localhost:8080/swagger/#/userdatas/userdatas%23setUserLikeArt (POST)
  
  

for e.g:
<pre>
curl -X 'POST' \
'http://localhost:8080/userdatas/create' \
-H 'accept: application/json' \
-H 'Content-Type: multipart/form-data' \
-F 'twitter_link=https://www.twitter.com/@Williams_Scottish' \
-F 'realname=Williams Scottish' \
-F 'cover_photo={
"content": "VmVsIGFyY2hpdGVjdG8gcXVpYSBldCBkb2xvcmVzLg=="
}' \
-F 'artist_pastelid=jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS' \
-F 'facebook_link=https://www.facebook.com/Williams_Scottish' \
-F 'biography=I'\''m a digital artist based in Paris, France. ...' \
-F 'primary_language=English' \
-F 'avatar_image={
"content": "VmVsIGFyY2hpdGVjdG8gcXVpYSBldCBkb2xvcmVzLg=="
}' \
-F 'location=New York, US' \
-F 'native_currency=USD' \
-F 'categories=Quos dignissimos ut corrupti eos aut.' \
-F 'artist_pastelid_passphrase=qwerasdf1234'


curl -X 'POST' \
  'http://localhost:8080/userdatas/update' \
  -H 'accept: application/json' \
  -H 'Content-Type: multipart/form-data' \
  -F 'twitter_link=' \
  -F 'realname=' \
  -F 'cover_photo=' \
  -F 'artist_pastelid=jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VaoRqi1GnQrYKjSxQAC7NBtvtEdS' \
  -F 'facebook_link=' \
  -F 'biography=' \
  -F 'primary_language=' \
  -F 'avatar_image=' \
  -F 'location=' \
  -F 'native_currency=' \
  -F 'categories=' \
  -F 'artist_pastelid_passphrase=qwerasdf9876'

  
curl -X 'GET' \
'http://localhost:8080/userdatas/jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS' \
-H 'accept: application/json'


--- SET FOLLOW RELATIONSHIP ---
curl -X 'POST' \
  'http://localhost:8080/userdatas/follow' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "followee_pastel_id": "jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VaoRqi1GnQrYKjSxQAC7NBtvtEdS",
  "follower_pastel_id": "jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS"
}'

--- GET FOLLOWERS ---
curl -X 'GET' \
  'http://localhost:8080/userdatas/follow/followers?pastelid=jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS&limit=10&offset=0' \
  -H 'accept: application/json'

--- GET FOLLOWEES ---
curl -X 'GET' \
  'http://localhost:8080/userdatas/follow/followees?pastelid=jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS&limit=10&offset=0' \
  -H 'accept: application/json'

--- GET FRIEND ---
curl -X 'GET' \
  'http://localhost:8080/userdatas/follow/friends?pastelid=jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS&limit=10&offset=0' \
  -H 'accept: application/json'
  
--- SET ART LIKE ---
curl -X 'POST' \
  'http://localhost:8080/userdatas/like/art' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "art_pastel_id": "jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VaoRqi1GnQrYKjSxQAC7NBtvtEdS",
  "user_pastel_id": "jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS"
}'

--- GET ART LIKE ---
curl -X 'GET' \
  'http://localhost:8080/userdatas/like/art' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "art_id": "jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS",
  "limit": 10,
  "offset": 0
}'


</pre>
  
  
* If we want to test rqlite clustering feature, need to do as described in function getDatabaseNodes

  
## IMPORTANT NOTICE

- When walletnode UI create the userdata, it can call the api /userdatas/create directly

- When walletnode UI update the userdata, it must read the current userdata info by call the api /userdatas/<pastelid> (GET), then receive complete user info, and choose any field(s) to update, then POST back the complete user info. Otherwise, there'll be violation the userdata's hash (for userdata's integrity)

  
## Troubleshooting

- supernode shutdown grpc server after serving the 1st request but not raising error

>Try to delete folder "supernode" in the path in step 2

## To use SNs as a cluster of rqlite
- Just need to remove if condition in gonode/supernode/cmd/app.go in function getDatabaseNodes()
```
	if test := sys.GetStringEnv("DB_CLUSTER", ""); test != "" {
		nodeList, err := pastelClient.MasterNodesList(ctx)
		if err != nil {
			return nil, err
		}
		for _, nodeInfo := range nodeList {
			address := nodeInfo.ExtAddress
			segments := strings.Split(address, ":")
			if len(segments) != 2 {
				return nil, errors.Errorf("malformed db node address: %s", address)
			}
			nodeAddress := fmt.Sprintf("%s:%d", segments[0], rqliteDefaultPort)
			nodeIPList = append(nodeIPList, nodeAddress)
		}
	}
```
- Or just need to make env DB_CLUSTER so that it has a value
