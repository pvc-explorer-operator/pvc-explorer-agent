# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /agent-bin ./cmd/agent

FROM gcr.io/distroless/static-debian12
COPY --from=builder /agent-bin /agent
ENTRYPOINT ["/agent"]
