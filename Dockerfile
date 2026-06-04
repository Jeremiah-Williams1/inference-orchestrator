# Multi-stage build.
# Stage 1 compiles the binary. Stage 2 runs it.
# The final image has no Go toolchain in it — just the binary.
# This keeps the image small and reduces the attack surface.

FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

# ---

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
