# Build the image
FROM golang:1.21-alpine AS Builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run Image
FROM golang:1.21-alpine
WORKDIR /appa
COPY --from=Builder /app/main .
COPY .env .
EXPOSE 8888
CMD  ["/app/main"]