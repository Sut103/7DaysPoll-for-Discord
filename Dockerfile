FROM golang:1.21 as build
WORKDIR /lambda
COPY go.mod go.sum ./
COPY lambda/main.go .
RUN GOARCH=amd64 GOOS=linux go build -tags 7dayspoll-build -o main main.go
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /lambda/main /main
ENTRYPOINT [ "/main" ]