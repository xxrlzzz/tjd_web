FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/traffic_jam_direction
COPY . $GOPATH/src/traffic_jam_direction
RUN go build .

EXPOSE 8000
ENTRYPOINT ["./traffic_jam_direction"]
