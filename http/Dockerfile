FROM shared-image:latest

RUN mkdir /http
ADD . /http/
WORKDIR /http
RUN go build -o main .
