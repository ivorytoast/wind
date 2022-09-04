FROM golang:alpine

LABEL maintainer="Anthony Hamill"

WORKDIR /opt/hamilla

COPY go.mod go.sum ./
RUN go mod download

COPY ./ .
RUN go build main.go

EXPOSE 10000
EXPOSE 8080
EXPOSE 5556

CMD ["./main"]
