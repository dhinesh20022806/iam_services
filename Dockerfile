FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server .

FROM alpine:3.21 AS final

RUN adduser -D appuser
USER appuser

WORKDIR /bin

COPY --from=buidler /app/server .

EXPOSE 50050

CMD ["./server"]