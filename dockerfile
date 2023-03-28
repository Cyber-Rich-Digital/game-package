FROM golang:1.20-alpine

WORKDIR /app

COPY . .
RUN apk update
RUN apk add alpine-sdk
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init
RUN go build -o ./build/API

EXPOSE 3000

CMD [ "./build/API" ]