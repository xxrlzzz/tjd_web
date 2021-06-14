FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/tjd_web
COPY . $GOPATH/src/tjd_web
#RUN go build .
RUN make build

EXPOSE 8000
ENTRYPOINT [".build/tjd_web"]
