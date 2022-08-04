FROM golang:1.18 as builder
LABEL maintainer="Yamac Eren Ay <yamac.ay@scandio.de>"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main
EXPOSE 8080
CMD ["./main"]