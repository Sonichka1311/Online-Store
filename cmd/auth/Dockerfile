FROM shared-image:latest

RUN mkdir auth
COPY . auth
WORKDIR auth
RUN go build -mod=mod -o main .

CMD ["./main"]