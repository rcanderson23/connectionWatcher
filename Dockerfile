FROM golang:1.16 as builder

ENV GOOS=linux
ENV CGO_ENABLED=0


RUN apt update && apt install git mercurial gcc -y
WORKDIR /tmp/app
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN go build -ldflags='-w -extldflags "-static"' -o /connectionWatcher

# Application
FROM alpine
RUN apk add iptables
COPY --from=builder /connectionWatcher /connectionWatcher
ENTRYPOINT ["/connectionWatcher"]
