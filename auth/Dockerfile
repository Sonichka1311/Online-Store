FROM shared-image:latest

RUN mkdir /auth
ADD . /auth/
WORKDIR /auth
RUN go build -o main .
