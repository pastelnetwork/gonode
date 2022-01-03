FROM debian

RUN apt-get update && apt-get install -y ca-certificates

COPY ./test_images /root/pastel
COPY ./p2pkeys.json /p2pkeys.json
COPY helper /helper


ENTRYPOINT [ "/helper" ]
