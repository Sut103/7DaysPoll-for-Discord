FROM golang:1.21 as build
WORKDIR /7dayspoll
COPY /app  .
RUN GOARCH=arm64 GOOS=linux go build -tags 7dayspoll-build -o main main.go
FROM gcr.io/distroless/base-debian12:latest-arm64
COPY --from=build /7dayspoll/main /main
ENTRYPOINT [ "/main" ]
