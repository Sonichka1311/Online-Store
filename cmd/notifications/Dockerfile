FROM shared-image:latest

RUN mkdir notifications
COPY . notifications
WORKDIR notifications
RUN go build -mod=mod -o main .

CMD ["./main"]