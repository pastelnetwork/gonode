FROM debian

RUN apt-get update && apt-get install -y ca-certificates

RUN mkdir -p /root/pastel

COPY helper /helper

ENTRYPOINT [ "/helper" ]
