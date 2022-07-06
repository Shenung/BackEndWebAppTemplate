FROM golang:latest

RUN mkdir /build
COPY . /build

WORKDIR /build

RUN go build ./backend/main.go
CMD ["./main"]



EXPOSE 8080