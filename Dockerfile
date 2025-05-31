FROM golang:1.24.3 AS build
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
RUN GOARCH=arm64 GOOS=linux go build -tags 7dayspoll-build -o main ./cmd/7dayspoll/main.go

FROM gcr.io/distroless/base-debian12:latest-arm64
COPY --from=build /app/main /main
ENTRYPOINT [ "/main" ]