FROM golang:1.23.4-alpine AS build-stage
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

# run tests now
FROM build-stage AS run-test-stage
RUN go test -v ./... -count=1

EXPOSE 23456
ENTRYPOINT ["/main"]
