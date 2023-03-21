FROM golang:1.19-alpine

WORKDIR /app/gophermart
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN go build -o ./gophermart ./cmd/gophermart

EXPOSE 8080 8080
CMD [ "/app/gophermart/gophermart" ]