FROM golang:1.23.0-alpine
WORKDIR /app/myapp
COPY go.mod go.sum ./
RUN go mod download
COPY ./myapp .
RUN go build -o main .
EXPOSE 8000
CMD ["./main"]