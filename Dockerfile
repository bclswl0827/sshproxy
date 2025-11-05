FROM golang:alpine AS builder
# Uncomment the following line to use a mirror of go module proxy
# ENV GOPROXY="https://goproxy.cn,direct"
COPY . /build_src
WORKDIR /build_src
RUN go build -v

FROM scratch
COPY --from=builder /build_src/sshproxy /sshproxy
ENTRYPOINT [ "/sshproxy" ]
