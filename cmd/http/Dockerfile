FROM shared-image:latest

RUN mkdir http
COPY . http
WORKDIR http
RUN go build -mod=vendor -o main .

CMD ["./main"]
