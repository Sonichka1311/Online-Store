FROM shared-image:latest

RUN mkdir uploader
COPY . uploader
WORKDIR uploader
RUN go build -mod=mod -o main .

CMD ["./main"]