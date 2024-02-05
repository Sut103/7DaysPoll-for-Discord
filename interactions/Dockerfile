FROM golang:1.21 as build
WORKDIR /interactions
COPY .  .
RUN GOARCH=amd64 GOOS=linux go build -tags 7dayspoll-build -o main main.go
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /interactions/main /main
ENTRYPOINT [ "/main" ]