FROM golang:1.21 as builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build

FROM debian:12-slim
LABEL org.opencontainers.image.source https://github.com/ondrejsika/counter
COPY --from=builder /build/counter /usr/local/bin/counter
CMD [ "/usr/local/bin/counter" ]
EXPOSE 80
