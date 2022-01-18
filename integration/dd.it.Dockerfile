FROM golang:1.17-alpine AS builder
RUN apk --update add ca-certificates git make g++
ENV GO111MODULE=on

COPY common/  /common/
COPY fakes/dd-server /fakes/dd-server/
COPY fakes/common /fakes/common/

WORKDIR /fakes/dd-server
RUN go mod download
RUN go build -o app

FROM golang:1.17-alpine
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /fakes/dd-server/app .

EXPOSE 50052
EXPOSE 51052

CMD [ "./app" ]
