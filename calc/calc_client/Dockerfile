FROM golang

RUN mkdir /app
WORKDIR /app
ADD . /app/client-service

ENTRYPOINT ["./client-service"]

EXPOSE 4000