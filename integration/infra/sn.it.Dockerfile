FROM golang:1.17-alpine AS builder
RUN apk --update add ca-certificates git make g++
ENV GO111MODULE=on

COPY supernode/ /supernode/
COPY hermes/ /hermes/
COPY common/  /common/
COPY dupedetection/ /dupedetection/
COPY go-webp/ /go-webp/
COPY mixins/ /mixins/
COPY bridge/ /bridge/
COPY metadb/ /metadb/
COPY pastel/ /pastel/
COPY p2p/ /p2p/
COPY proto/ /proto/
COPY raptorq/ /raptorq/
COPY /integration/configs/sn1.yml /supernode/sn1.yml
COPY /integration/configs/sn2.yml /supernode/sn2.yml
COPY /integration/configs/sn3.yml /supernode/sn3.yml
COPY /integration/configs/sn4.yml /supernode/sn4.yml
COPY /integration/configs/p1.conf /supernode/p1.conf
COPY /integration/configs/p2.conf /supernode/p2.conf
COPY /integration/configs/p3.conf /supernode/p3.conf
COPY /integration/configs/p4.conf /supernode/p4.conf
COPY integration/testdata/ /supernode/testdata/

WORKDIR /supernode
RUN go mod download
RUN go build -o app

FROM golang:1.17-alpine
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /supernode/app .
RUN mkdir configs
RUN mkdir testdata
COPY --from=builder /supernode/sn1.yml /configs/sn1.yml
COPY --from=builder /supernode/sn2.yml /configs/sn2.yml
COPY --from=builder /supernode/sn3.yml /configs/sn3.yml
COPY --from=builder /supernode/sn4.yml /configs/sn4.yml
COPY --from=builder /supernode/p1.conf /configs/p1.conf
COPY --from=builder /supernode/p2.conf /configs/p2.conf
COPY --from=builder /supernode/p3.conf /configs/p3.conf
COPY --from=builder /supernode/p4.conf /configs/p4.conf
COPY --from=builder /supernode/testdata/ /testdata/



EXPOSE 9090
EXPOSE 14444
CMD [ "./app" ]
