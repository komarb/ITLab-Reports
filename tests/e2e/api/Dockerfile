FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY ./itlab-reports-api-build/ .

CMD ./main
