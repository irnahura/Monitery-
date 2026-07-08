FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN go build -o /out/backend ./cmd/api

FROM alpine:3.21
WORKDIR /app
RUN adduser -D -H appuser
COPY --from=build /out/backend /app/backend
USER appuser
EXPOSE 8080
CMD ["/app/backend"]
