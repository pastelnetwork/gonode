FROM golang:1.17-alpine AS builder
RUN apk --update add ca-certificates git make g++
ENV GO111MODULE=on

COPY supernode/ /supernode/
COPY common/  /common/
COPY dupedetection/ /dupedetection/
COPY go-webp/ /go-webp/
COPY metadb/ /metadb/
COPY pastel/ /pastel/
COPY p2p/ /p2p/
COPY proto/ /proto/
COPY raptorq/ /raptorq/
COPY /integration/configs/sn1.yml /supernode/sn1.yml
COPY /integration/configs/sn2.yml /supernode/sn2.yml
COPY /integration/configs/sn3.yml /supernode/sn3.yml
COPY /integration/configs/p1.conf /supernode/p1.conf
COPY /integration/configs/p2.conf /supernode/p2.conf
COPY /integration/configs/p3.conf /supernode/p3.conf

WORKDIR /supernode
RUN go mod download
RUN go build -o app

FROM golang:1.17-alpine
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /supernode/app .
RUN mkdir configs
COPY --from=builder /supernode/sn1.yml /configs/sn1.yml
COPY --from=builder /supernode/sn2.yml /configs/sn2.yml
COPY --from=builder /supernode/sn3.yml /configs/sn3.yml
COPY --from=builder /supernode/p1.conf /configs/p1.conf
COPY --from=builder /supernode/p2.conf /configs/p2.conf
COPY --from=builder /supernode/p3.conf /configs/p3.conf


EXPOSE 9090
CMD [ "./app" ]
