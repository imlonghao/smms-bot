FROM golang:1.12-alpine AS builder
LABEL maintainer="imlonghao <dockerfile@esd.cc>"
WORKDIR /builder
COPY . /builder
RUN apk add upx && \
    GO111MODULE=on go build -mod=vendor -ldflags="-s -w" -o /app && \
    upx --lzma --best /app

FROM gcr.io/distroless/base
COPY --from=builder /app .
CMD ["/app"]