FROM golang:1.12

WORKDIR /app

RUN go get github.com/gorilla/mux
RUN go get go.mongodb.org/mongo-driver/bson
RUN go get go.mongodb.org/mongo-driver/mongo

COPY . .

CMD ["go","run","main.go"]
