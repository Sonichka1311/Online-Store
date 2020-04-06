FROM shared-image:latest

RUN mkdir mock
COPY . mock
WORKDIR mock
RUN go build -mod=vendor -o main .

CMD ["./main"]