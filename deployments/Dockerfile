# Etapa de construcción
FROM golang:1.26-alpine AS build

WORKDIR /app

COPY . .
RUN go mod tidy

RUN go build -o app ./cmd/api

# Etapa de ejecución
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=build /app/app .

EXPOSE 8080

CMD ["./app"]
