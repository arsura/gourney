FROM golang:1.16-alpine AS builder

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /app

COPY . .

RUN go build -o main .

FROM alpine

COPY --from=builder /app/main .

CMD ["./main"]