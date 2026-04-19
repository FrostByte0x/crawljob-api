# Stage 1
FROM golang:1.26.2 AS builder

WORKDIR /app
COPY . . 
ENV GOOS=linux
ENV GOARCH=amd64

RUN CGO_ENABLED=0 go build -o crawljob-api .

# Stage 2
FROM alpine

COPY --from=builder /app/crawljob-api /crawljob-api

ENV DESTINATION_FOLDER=/mnt/jDownloader/crawljob-api
ENV CRAWLJOB_FOLDER=/mnt/jDownloader/crawljobs
ENV ALLOWED_DOMAINS=1fichier.com,mega.nz
ENV ENABLE_PURGE=false

CMD ["/crawljob-api"]

EXPOSE 8080
