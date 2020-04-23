FROM shared-image:latest

RUN mkdir http
COPY . http
WORKDIR http
RUN go build -mod=mod -o main .

CMD ["./main"]
