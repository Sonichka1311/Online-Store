FROM golang:latest

RUN mkdir /common
RUN mkdir /utils

ADD ./common /common/
ADD ./utils /utils/

RUN go get github.com/gorilla/mux
RUN go get github.com/dgrijalva/jwt-go