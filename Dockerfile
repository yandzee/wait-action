FROM golang:1.21-alpine AS preparation

COPY . ./src
RUN go build -o /action ./src/main.go

ENTRYPOINT ["/action"]
