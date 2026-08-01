FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /weather-api ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /weather-api /weather-api

EXPOSE 8080

ENTRYPOINT ["/weather-api"]
