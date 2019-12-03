FROM golang:1.12.4-alpine3.9 as builder

WORKDIR /go/src/heartbeat_demo

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o heartbeat_demo .

FROM alpine:3.9 as prod

WORKDIR /root/heartbeat_demo

COPY --from=0 /go/src/heartbeat_demo .

#EXPOSE 8080


CMD ["./heartbeat_demo"]
