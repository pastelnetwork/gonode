FROM golang:1.17-alpine AS builder
RUN apk --update add ca-certificates git make g++
ENV GO111MODULE=on

COPY walletnode/ /walletnode/
COPY common/  /common/
COPY dupedetection/ /dupedetection/
COPY go-webp/ /go-webp/
COPY metadb/ /metadb/
COPY pastel/ /pastel/
COPY p2p/ /p2p/
COPY proto/ /proto/
COPY raptorq/ /raptorq/
COPY /integration/configs/wn.yml /walletnode/wn.yml
COPY /integration/configs/p4.conf /walletnode/p4.conf

WORKDIR /walletnode
RUN go mod download
RUN go build -o app

FROM golang:1.17-alpine
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /walletnode/app .
RUN mkdir configs
COPY --from=builder /walletnode/wn.yml /configs/wn.yml
COPY --from=builder /walletnode/p4.conf /configs/p4.conf

EXPOSE 8089
CMD [ "./app" ]
