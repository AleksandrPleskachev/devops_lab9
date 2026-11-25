FROM golang:1.22.12-alpine AS builder
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/service .
FROM alpine:3.19
RUN apk add --no-cache tzdata && adduser -D app
COPY --from=builder /out/service /service
USER app
EXPOSE 8080
HEALTHCHECK --interval=3s --timeout=2s --retries=2 CMD ["/service", "probe"]
ENTRYPOINT ["/service"]
