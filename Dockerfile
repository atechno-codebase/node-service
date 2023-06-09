FROM golang:1.20.4-alpine AS build

WORKDIR /build

COPY . .

RUN go mod tidy

RUN go build -o app

FROM alpine:latest

WORKDIR /app

COPY --from=build /build/app /app/app

EXPOSE 8000

USER nonroot:nonroot

ENTRYPOINT ["/app"]
