FROM shared-image:latest

RUN mkdir auth
WORKDIR auth
COPY . .
RUN go build -mod=mod -o main .

CMD ["./main"]