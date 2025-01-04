FROM golang:1.21.6 AS builder
RUN git clone https://github.com/magefile/mage && \
    cd mage && \
    go run bootstrap.go

ADD . /go/src/infra-example
WORKDIR /go/src/infra-example
RUN mage services

FROM alpine:3
RUN apk add --no-cache libc6-compat
WORKDIR /root/infra-example/
COPY --from=builder /go/src/infra-example/build/* ./
ENV PATH="/root/infra-example:${PATH}"
