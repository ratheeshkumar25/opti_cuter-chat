FROM golang:1.23.1-alpine AS builder 

WORKDIR /app

COPY . /app

RUN go build -o chatservice ./cmd

FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/chatservice .

COPY .env /app/

CMD ["./chatservice"]