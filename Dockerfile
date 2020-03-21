#FROM tarantool/tarantool:2.2.1
#COPY databaseServer.lua /opt/tarantool
#CMD ["tarantool", "/opt/tarantool/databaseServer.lua"]

FROM golang:latest

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -o main .
#RUN chmod +x run.sh
#ENTRYPOINT ["/app/run.sh"]
