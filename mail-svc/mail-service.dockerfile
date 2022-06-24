FROM alpine:latest

RUN mkdir /app

COPY ./bin/mailerApp /app
COPY ./pkg/lib/templates /templates

CMD ["/app/mailerApp"]
