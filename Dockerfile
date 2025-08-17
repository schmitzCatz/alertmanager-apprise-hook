FROM golang:1.25.0 AS builder

WORKDIR /app
COPY ../go.mod go.sum *.go ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /apprise-webhook


FROM alpine:3.22.1
COPY --from=builder /apprise-webhook /app

EXPOSE 8080

ENV TAG="all"
ENV APPRISE_URL=""
ENV LISTEN_ADDRESS=":8080"

CMD ["/app"]

LABEL org.opencontainers.image.authors="Oliver Schmitz"
LABEL org.opencontainers.image.url="https://github.com/schmitzCatz/alertmanager-apprise-hook"
LABEL org.opencontainers.image.documentation="https://github.com/schmitzCatz/alertmanager-apprise-hook"
LABEL org.opencontainers.image.source="https://github.com/schmitzCatz/alertmanager-apprise-hook"
LABEL org.opencontainers.image.vendor="Oliver Schmitz"
LABEL org.opencontainers.image.licenses="GNU GPLv3"
LABEL org.opencontainers.image.title="Alertmanager Apprise Webhook"
LABEL org.opencontainers.image.description="A webwook for Prometheus Alertmanager to forwars notifications to Apprise"