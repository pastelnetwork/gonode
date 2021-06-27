This file for sample reference testing only, should be remove at PR

Setup raptorQ service:
- install rust : nurl --proto '=https' --tlsv1.2 https://sh.rustup.rs -sSf | sh
- clone at https://github.com/pastelnetwork/rqservice
- build service : cargo build
- run service : ./target/debug/rq-service  -s 127.0.0.1:50051

