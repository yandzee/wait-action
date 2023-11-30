FROM golang:1.21-alpine AS preparation

WORKDIR /app
COPY . .
RUN go build -o /bin/wait-action main.go

# RUN cp ./wait-action $GOPATH/
# WORKDIR $GOPATH

# CMD ./wait-action
ENTRYPOINT ["/bin/wait-action"]
