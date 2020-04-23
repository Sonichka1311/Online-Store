FROM golang:latest

RUN mkdir src/shop
RUN mkdir src/shop/pkg
COPY pkg src/shop/pkg
COPY go.mod src/shop
WORKDIR src/shop

RUN go mod download