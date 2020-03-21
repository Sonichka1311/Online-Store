FROM tarantool/tarantool:1.10.2

COPY databaseServer.lua /opt/tarantool
WORKDIR /opt/tarantool
RUN mkdir ./logs
RUN touch ./logs/logFile.txt
RUN chmod 777 -R /opt/tarantool
CMD ["tarantool", "databaseServer.lua"]
