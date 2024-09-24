FROM golang:latest

WORKDIR /src
ADD . /src

RUN go install github.com/air-verse/air@latest

CMD ["air", "-c", ".air.toml"]