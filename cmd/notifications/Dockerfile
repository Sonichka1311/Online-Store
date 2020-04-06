FROM shared-image:latest

RUN mkdir notifications
COPY . notifications
WORKDIR notifications
RUN go build -mod=vendor -o main .

CMD ["./main"]