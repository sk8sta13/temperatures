FROM golang:1.24 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /apitemperature cmd/api/main.go

FROM scratch
WORKDIR /
COPY --from=build /apitemperature /apitemperature
ENTRYPOINT ["/apitemperature"]