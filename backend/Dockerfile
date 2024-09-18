FROM golang:latest as builder
WORKDIR /src
ADD . /src/

RUN CGO_ENABLED=0 GOOS=linux go build -o build .

FROM scratch as prod
COPY --from=builder /src/build /
CMD ["./build"]