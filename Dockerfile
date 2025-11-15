FROM golang:1.25.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o pr_reviewer ./cmd/pr_reviewer/main.go

FROM alpine:latest AS app
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/pr_reviewer /app/pr_reviewer
EXPOSE 8080

CMD ["./pr_reviewer"]
