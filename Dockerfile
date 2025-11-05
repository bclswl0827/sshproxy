FROM golang:alpine AS builder
# Uncomment the following line to use a mirror of APK repository
ENV APK_SOURCE_HOST="mirrors.bfsu.edu.cn"
# Uncomment the following line to use a mirror of go module proxy
ENV GOPROXY="https://goproxy.cn,direct"
COPY . /build_src
WORKDIR /build_src
RUN if [ "x${APK_SOURCE_HOST}" != "x" ]; then \
    sed -i "s/dl-cdn.alpinelinux.org/$APK_SOURCE_HOST/g" /etc/apk/repositories; \
    fi \
    && apk add --update --no-cache ca-certificates\
    && go build -v -o sshproxy .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build_src/sshproxy /sshproxy
ENTRYPOINT ["/sshproxy"]
