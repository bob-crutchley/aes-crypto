FROM golang:1.11.5-alpine3.9 as build
WORKDIR /build
COPY . .
RUN go build *.go
FROM alpine:3.9
WORKDIR /app
COPY --from=build /build/aes-crypto .
EXPOSE 8000
ENTRYPOINT ["./aes-crypto"]
