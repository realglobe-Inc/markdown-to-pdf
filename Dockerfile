FROM golang:1.13 as builder
WORKDIR /app
#COPY go.* ./
#RUN go mod download
COPY invoke.go ./
#RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o server
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server

FROM alpine:3
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/server /server
COPY script.sh ./
CMD ["/server"]
