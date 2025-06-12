FROM golang:1.22.3-alpine3.19 AS build

RUN apk --no-cache add gcc g++ make ca-certificates

WORKDIR /go/src/github.com/eyepatch5263/go-grpc-microservices

COPY go.mod go.sum ./

COPY vendor vendor

COPY product product

RUN go build -mod vendor -o /go/bin/app ./product/cmd/product

FROM alpine:3.19

WORKDIR /usr/bin

COPY --from=build /go/bin/app .

EXPOSE 8080

CMD ["app"]