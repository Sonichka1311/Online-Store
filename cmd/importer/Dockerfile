FROM shared-image:latest

RUN mkdir importer
COPY . importer
WORKDIR importer
RUN go build -mod=mod -o main .

CMD ["./main"]