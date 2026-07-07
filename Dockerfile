FROM golang:1.24.3 AS build
WORKDIR /7dayspoll
COPY /app  .
RUN GOARCH=amd64 GOOS=linux go build -tags 7dayspoll-build -o main main.go
FROM gcr.io/distroless/base-debian12:latest
COPY --from=build /7dayspoll/main /main
ENTRYPOINT [ "/main" ]
