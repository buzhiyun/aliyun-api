FROM golang:1.25-alpine as build

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn

ADD . /app
WORKDIR /app


RUN  sed -i 's#dl-cdn.alpinelinux.org#mirrors.cloud.tencent.com#g' /etc/apk/repositories && \
     go mod vendor && sed -i '/https:\/\/fonts.googleapis.com/d' vendor/github.com/iris-contrib/swagger/v12/swagger.go && \
     go build -ldflags '-s -w' -o aliyun-api aliyun.go



FROM alpine:3.23
RUN sed -i 's#dl-cdn.alpinelinux.org#mirrors.aliyun.com#g' /etc/apk/repositories  && apk add sudo curl && \
    sed -i 's#mirrors.aliyun.com#mirrors.cloud.aliyuncs.com#g' /etc/apk/repositories  && \
    rm -rf /var/cache/apk/* && \
    rm -rf /root/.cache && \
    rm -rf /tmp/* && \
    echo 'lucifer ALL=(ALL) ALL,NOPASSWD:/sbin/apk' >> /etc/sudoers && \
    adduser -h /app -u 1000 -D lucifer && \
    echo -e '\n\n# septnet CA' >> /etc/ssl/certs/ca-certificates.crt && curl 'https://7netpublic.oss-cn-hangzhou.aliyuncs.com/dev/ca/septnet-ca.crt' >> /etc/ssl/certs/ca-certificates.crt
    

# 不要用root
USER lucifer
WORKDIR /app

COPY --from=build /app/aliyun-api .

CMD ["/app/aliyun-api"]
