FROM golang:1.21.6 AS builder
RUN git clone https://github.com/magefile/mage && \
    cd mage && \
    go run bootstrap.go

ADD . /go/src/infra-example
WORKDIR /go/src/infra-example
RUN mage services
RUN mage -compile ./build/mage

FROM alpine:3
RUN apk add --no-cache libc6-compat
WORKDIR /root/infra-example/
COPY --from=builder /go/src/infra-example/build/* ./
COPY --from=builder /go/src/infra-example/docs ./docs/
ENV PATH="/root/infra-example:${PATH}"
