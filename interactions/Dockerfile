FROM golang:1.21 as build
WORKDIR /interactions
COPY .  .
RUN GOARCH=arm64 GOOS=linux go build -tags 7dayspoll-build -o main main.go
FROM gcr.io/distroless/base-debian12:latest-arm64
COPY --from=build /interactions/main /main
ENTRYPOINT [ "/main" ]
