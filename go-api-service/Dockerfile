FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server ./ 

FROM alpine:3.18
COPY --from=build /app/server /usr/local/bin/server
EXPOSE 8081
ENTRYPOINT ["/usr/local/bin/server"]