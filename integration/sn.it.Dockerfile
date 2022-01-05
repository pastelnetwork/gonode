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
COPY /integration/configs/pastel.conf /supernode/pastel.conf

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
COPY --from=builder /supernode/pastel.conf /configs/pastel.conf

EXPOSE 9090
CMD [ "./app" ]
