FROM golang:1.12

WORKDIR /app

RUN go get github.com/gorilla/mux
RUN go get go.mongodb.org/mongo-driver/bson
RUN go get go.mongodb.org/mongo-driver/mongo
RUN go get github.com/auth0/go-jwt-middleware
RUN go get github.com/dgrijalva/jwt-go

COPY . .

CMD ["go","run","main.go"]
