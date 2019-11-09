FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY ./itlab-reports-api-build/ .
RUN chmod +x main

COPY ./config.testmace.json/ ./config.json
CMD ./main
