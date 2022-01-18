FROM golang:1.17-alpine AS builder
RUN apk --update add ca-certificates git make g++
ENV GO111MODULE=on

COPY common/  /common/
COPY fakes/pasteld /fakes/pasteld/
COPY fakes/common /fakes/common/

WORKDIR /fakes/pasteld
RUN go mod download
RUN go build -o app

FROM golang:1.17-alpine
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /fakes/pasteld/app .

EXPOSE 19932
EXPOSE 29932

CMD [ "./app" ]