FROM golang:1.23
WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

ADD go.mod .
ADD go.sum .

RUN go mod download
COPY . .

RUN go build -o ./bot
CMD ["/app/bot"]