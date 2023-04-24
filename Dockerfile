FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY .env .
COPY google_tags.txt .
COPY --from=builder /app/app .

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

CMD ["./app"]
