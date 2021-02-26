FROM golang

RUN mkdir /app
WORKDIR /app
ADD . /app/server-service

ENTRYPOINT ["./server-service"]

EXPOSE 4000