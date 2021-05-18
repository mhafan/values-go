FROM golang:latest

WORKDIR /go/src

RUN git clone https://github.com/mhafan/values-go.git

WORKDIR /go/src/values-go

RUN bash compile

CMD cnt/cnt -h my_redis_web:6379 -X
