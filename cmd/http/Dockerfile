FROM shared-image:latest

RUN mkdir /http
COPY . /http
WORKDIR /http
RUN go build -o main .

CMD ["./main"]
