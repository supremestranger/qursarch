FROM golang:latest as builder
WORKDIR /src
ADD . /src/

RUN CGO_ENABLED=0 go build -ldflags=all="-w -s" -o build .

FROM scratch as dev
EXPOSE 3001
COPY --from=builder /src/build /
CMD ["./build"]