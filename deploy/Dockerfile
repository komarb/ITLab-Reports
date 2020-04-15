FROM alpine

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

WORKDIR /app

COPY ITLab-Reports/itlab-reports-api-build/main ./main

RUN chmod +x main

CMD [ "./main" ]