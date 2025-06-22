FROM golang:1.20-alpine

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o match-service

EXPOSE 8083

CMD ["./match-service"]
