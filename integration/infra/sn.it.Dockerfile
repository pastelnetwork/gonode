FROM golang:1.17 AS builder
RUN apt-get update 
ENV GO111MODULE=on

COPY supernode/ /supernode/
COPY hermes/ /hermes/
COPY common/  /common/
COPY dupedetection/ /dupedetection/
COPY go-webp/ /go-webp/
COPY mixins/ /mixins/
COPY metadb/ /metadb/
COPY pastel/ /pastel/
COPY p2p/ /p2p/
COPY proto/ /proto/
COPY raptorq/ /raptorq/
COPY /integration/configs/sn1.yml /supernode/sn1.yml
COPY /integration/configs/sn2.yml /supernode/sn2.yml
COPY /integration/configs/sn3.yml /supernode/sn3.yml
COPY /integration/configs/sn4.yml /supernode/sn4.yml
COPY /integration/configs/sn5.yml /supernode/sn5.yml
COPY /integration/configs/sn6.yml /supernode/sn6.yml
COPY /integration/configs/sn7.yml /supernode/sn7.yml
COPY /integration/configs/p1.conf /supernode/p1.conf
COPY /integration/configs/p2.conf /supernode/p2.conf
COPY /integration/configs/p3.conf /supernode/p3.conf
COPY /integration/configs/p4.conf /supernode/p4.conf
COPY /integration/configs/p5.conf /supernode/p5.conf
COPY /integration/configs/p6.conf /supernode/p6.conf
COPY /integration/configs/p7.conf /supernode/p7.conf
COPY /integration/configs/rqservice.toml /supernode/rqservice.toml
COPY integration/testdata/ /supernode/testdata/
RUN curl https://download.pastel.network/beta/pastelup/pastelup-linux-amd64 -o /supernode/pastelup

WORKDIR /supernode
RUN go mod download
RUN go build -o app

FROM golang:1.17
RUN mkdir .pastel
COPY --from=builder /supernode/p1.conf /root/.pastel/pastel.conf
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /supernode/app .
COPY --from=builder /supernode/pastelup .
RUN chmod a+x pastelup

RUN mkdir configs
RUN mkdir testdata
RUN ./pastelup install rq-service -r beta
COPY --from=builder /supernode/rqservice.toml /root/.pastel/rqservice.toml
COPY --from=builder /supernode/sn1.yml /configs/sn1.yml
COPY --from=builder /supernode/sn2.yml /configs/sn2.yml
COPY --from=builder /supernode/sn3.yml /configs/sn3.yml
COPY --from=builder /supernode/sn4.yml /configs/sn4.yml
COPY --from=builder /supernode/sn5.yml /configs/sn5.yml
COPY --from=builder /supernode/sn6.yml /configs/sn6.yml
COPY --from=builder /supernode/sn7.yml /configs/sn7.yml
COPY --from=builder /supernode/p1.conf /configs/p1.conf
COPY --from=builder /supernode/p2.conf /configs/p2.conf
COPY --from=builder /supernode/p3.conf /configs/p3.conf
COPY --from=builder /supernode/p4.conf /configs/p4.conf
COPY --from=builder /supernode/p5.conf /configs/p5.conf
COPY --from=builder /supernode/p6.conf /configs/p6.conf
COPY --from=builder /supernode/p7.conf /configs/p7.conf
COPY --from=builder /supernode/testdata/ /testdata/



EXPOSE 9090
EXPOSE 50051
EXPOSE 14444
CMD [ "./app" ]
