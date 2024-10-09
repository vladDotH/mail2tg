FROM golang:1.23.2-alpine3.20 AS builder
WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

ADD go.mod .
ADD go.sum .

RUN go mod download
COPY . .

RUN go build -o ./bot

FROM alpine:3.20
COPY --from=builder /app /app
CMD ["/app/bot"]