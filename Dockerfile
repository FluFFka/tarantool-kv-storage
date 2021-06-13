FROM golang:alpine
COPY . /service
WORKDIR /service
RUN go build -o bin/main cmd/main.go
CMD ["/service/bin/main"]