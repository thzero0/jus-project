FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o api ./cmd

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/api .
EXPOSE 8080
CMD ["./api"]
