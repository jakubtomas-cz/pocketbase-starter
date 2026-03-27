FROM golang:1.25-alpine AS builder

WORKDIR /app
ADD . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/server/main.go

FROM alpine:3.21 AS production

WORKDIR /app
COPY --from=builder /app/app ./app
COPY --from=builder /app/views ./views
COPY --from=builder /app/pb_public ./pb_public

EXPOSE 1080
VOLUME ["/app/pb_data"]
CMD ["./app", "serve", "--http=0.0.0.0:8090"]
