# build
FROM golang:1.24.3 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/web

# run
FROM gcr.io/distroless/base-debian12
COPY --from=build /bin/app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
