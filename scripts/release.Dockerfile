FROM golang:1.17 AS builder
RUN apt-get update && \
    apt-get -y install gcc mono-mcs && \
    rm -rf /var/lib/apt/lists/*
RUN apt-get update
RUN apt-get install -y gcc-mingw-w64-x86-64 
RUN apt-get install -y gcc-mingw-w64-i686 
ENV GO111MODULE=on
#ARG LD_FLAGS
ENV LD_FLAGS ${LD_FLAGS}

COPY walletnode/ /walletnode/
COPY common/  /common/
COPY dupedetection/ /dupedetection/
COPY go-webp/ /go-webp/
COPY metadb/ /metadb/
COPY pastel/ /pastel/
COPY p2p/ /p2p/
COPY proto/ /proto/
COPY raptorq/ /raptorq/
COPY mixins/ /mixins/
COPY supernode/ /supernode/
COPY hermes/ /hermes/
COPY mixins/ /mixins/

WORKDIR /walletnode
RUN go mod download
RUN CGO_ENABLED=1 GOOS=windows GOARCH=amd64  CC=x86_64-w64-mingw32-gcc go build -ldflags=${LD_FLAGS} -o  walletnode-win64.exe
RUN CGO_ENABLED=1 GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc  go build -ldflags=${LD_FLAGS} -o walletnode-win32.exe
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags=${LD_FLAGS} -o walletnode-linux-amd64

WORKDIR /supernode
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags=${LD_FLAGS} -o supernode-linux-amd64