FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN apk add build-base && go build -o forum ./cmd/web/*
FROM alpine:3.6
LABEL Authors="@nzharylk && @tlsh0 && @abaltaba" Project="Forum"
WORKDIR /app
COPY --from=builder /app .
CMD ["/app/forum"]