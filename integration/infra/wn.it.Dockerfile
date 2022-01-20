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
COPY integration/testdata/ /walletnode/testdata/
COPY /integration/configs/wn.yml /walletnode/wn.yml
COPY /integration/configs/p0.conf /walletnode/p0.conf


WORKDIR /walletnode
RUN go mod download
RUN go build -o app

FROM golang:1.17-alpine
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /walletnode/app .
RUN mkdir configs
RUN mkdir testdata
COPY --from=builder /walletnode/wn.yml /configs/wn.yml
COPY --from=builder /walletnode/p0.conf /configs/p0.conf
COPY --from=builder /walletnode/testdata/09f6c459-ec2a-4db1-a8fe-0648fd97b5cb /testdata/09f6c459-ec2a-4db1-a8fe-0648fd97b5cb

EXPOSE 8089
CMD [ "./app" ]
