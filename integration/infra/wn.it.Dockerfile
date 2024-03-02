FROM golang:1.17 AS builder
RUN apt-get update 
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
COPY mixins/ /mixins/
COPY integration/testdata/ /walletnode/testdata/
COPY /integration/configs/wn.yml /walletnode/wn.yml
COPY /integration/configs/p0.conf /walletnode/p0.conf
COPY /integration/configs/rqservice.toml /walletnode/rqservice.toml
RUN curl https://download.pastel.network/beta/pastelup/pastelup-linux-amd64 -o /walletnode/pastelup


WORKDIR /walletnode
RUN go mod download
RUN go build -o app

FROM golang:1.17
RUN mkdir .pastel
COPY --from=builder /walletnode/p0.conf /root/.pastel/pastel.conf
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /walletnode/app .
COPY --from=builder /walletnode/pastelup .
RUN chmod a+x pastelup

RUN mkdir configs
RUN mkdir testdata
RUN ./pastelup install rq-service -r beta
COPY --from=builder /walletnode/wn.yml /configs/wn.yml
COPY --from=builder /walletnode/p0.conf /configs/p0.conf
COPY --from=builder /walletnode/rqservice.toml /root/.pastel/rqservice.toml
COPY --from=builder /walletnode/testdata/09f6c459-ec2a-4db1-a8fe-0648fd97b5cb /testdata/09f6c459-ec2a-4db1-a8fe-0648fd97b5cb

EXPOSE 8089
EXPOSE 50051
CMD [ "./app" ]
