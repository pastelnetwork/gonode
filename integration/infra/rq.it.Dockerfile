FROM golang:1.17
RUN apt-get update 
ENV GO111MODULE=on

RUN curl https://download.pastel.network/beta/pastelup-linux-amd64 -o pastelup
RUN chmod a+x pastelup
RUN mkdir .pastel
RUN ls -a 
RUN ./pastelup install rq-service -r beta
COPY /integration/configs/p0.conf /root/.pastel/pastel.conf
COPY /integration/configs/rqservice.toml /root/.pastel/rqservice.toml

EXPOSE 50051

CMD ["./pastelup", "start", "rq-service"]