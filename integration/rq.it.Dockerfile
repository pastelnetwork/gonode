FROM golang:1.17-alpine AS builder
RUN apk --update add ca-certificates git make g++
ENV GO111MODULE=on

COPY common/  /common/
COPY fakes/rq-service /fakes/rq-service/
COPY fakes/common /fakes/common/

WORKDIR /fakes/rq-service
RUN go mod download
RUN go build -o app

FROM golang:1.17-alpine
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /fakes/rq-service/app .

EXPOSE 50051
EXPOSE 51051

CMD [ "./app" ]