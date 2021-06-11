FROM golang:alpine
COPY . /service
WORKDIR /service
RUN apk add make && make build
CMD ["/service/bin/main"]