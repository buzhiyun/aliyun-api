FROM golang:1.18.10-alpine3.17 as build

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn

ADD . /app
WORKDIR /app


RUN  sed -i 's#dl-cdn.alpinelinux.org#mirrors.cloud.tencent.com#g' /etc/apk/repositories && \
     go mod vendor && go build -ldflags '-s -w' -o aliyun-api aliyun.go



FROM alpine:3.17
# 不要用root
RUN sed -i 's#dl-cdn.alpinelinux.org#mirrors.cloud.tencent.com#g' /etc/apk/repositories  && apk add sudo && \
    rm -rf /var/cache/apk/* && \
    rm -rf /root/.cache && \
    rm -rf /tmp/* && \
    echo 'lucifer ALL=(ALL) ALL,NOPASSWD:/sbin/apk' >> /etc/sudoers && \
    adduser -h /app -u 1000 -D lucifer

USER lucifer
WORKDIR /app

COPY --from=build /app/aliyun-api .

CMD ["/app/aliyun-api"]
