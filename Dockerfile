FROM golang:1.17.4
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY index.go ./main.go

RUN go build -o main .
EXPOSE 8001

CMD ["./main"]
