FROM golang:latest

RUN mkdir src/shop
RUN mkdir src/shop/pkg
COPY . src/shop/pkg

RUN go get github.com/gorilla/mux
RUN go get github.com/dgrijalva/jwt-go