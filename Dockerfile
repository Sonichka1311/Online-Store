FROM golang:latest

RUN mkdir src/shop
COPY go.mod src/shop
WORKDIR src/shop
RUN go mod download
RUN mkdir pkg
COPY pkg pkg