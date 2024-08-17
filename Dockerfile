# Build the go app.
FROM golang:1.22.2 AS build-stage

WORKDIR /backend

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /backend/backend/server
RUN go build -o server .

EXPOSE 5173
EXPOSE 8080

CMD ./server
