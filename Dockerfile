FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/capital-gains ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/capital-gains /app/capital-gains

ENTRYPOINT ["/app/capital-gains"]