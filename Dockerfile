FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest

ENV PATH="/root/go/bin:${PATH}"

RUN swag init -g ./cmd/main.go


RUN CGO_ENABLED=0 GOOS=linux go build -o task_manager ./cmd/main.go

# ---------------- Runtime ----------------

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root

COPY --from=builder /app/task_manager .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/.env .

EXPOSE 3000

CMD ["./task_manager"]