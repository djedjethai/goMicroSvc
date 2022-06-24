FROM alpine:latest

RUN mkdir /app

COPY ./bin /app

CMD ["/app/loggerApp"]
