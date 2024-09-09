FROM golang:1.23
WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

ADD go.mod .
ADD go.sum .

RUN go mod download
COPY --exclude=*.env . .

RUN go build -o ./bot
COPY *.env .

CMD ["/app/bot"]