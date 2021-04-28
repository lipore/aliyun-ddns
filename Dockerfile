FROM golang:1.16-alpine3.12 as builder

WORKDIR /opt/src/
ADD . .

RUN GOPROXY=https://goproxy.cn,direct go build && ls -al

FROM alpine:3.12

WORKDIR /opt/aliyun-ddns
COPY --from=builder /opt/src/aliyun-ddns .

CMD [ "/opt/aliyun-ddns/aliyun-ddns", "/opt/aliyun-ddns/config/config.yml" ]