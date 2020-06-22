FROM golang:1.14.2 as builder
WORKDIR /app
COPY invoke.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server

FROM debian:stable-20200607-slim
RUN apt-get update \
    && apt-get install -y pandoc texlive-xetex texlive-lang-japanese \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/server /server
COPY script.sh ./
CMD ["/server"]
