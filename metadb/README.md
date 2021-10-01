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

  
curl -X 'GET' \
'http://localhost:8080/userdatas/jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS' \
-H 'accept: application/json'

</pre>
  
  
* If we want to test rqlite clustering feature, need to do as described in function getDatabaseNodes

  
## IMPORTANT NOTICE

- When walletnode UI create the userdata, it can call the api /userdatas/create directly

- When walletnode UI update the userdata, it must read the current userdata info by call the api /userdatas/<pastelid> (GET), then receive complete user info, and choose any field(s) to update, then POST back the complete user info. Otherwise, there'll be violation the userdata's hash (for userdata's integrity)

  
## Troubleshooting

- supernode shutdown grpc server after serving the 1st request but not raising error

>Try to delete folder "supernode" in the path in step 2
