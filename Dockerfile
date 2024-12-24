FROM golang:1.23.4-alpine AS build-stage
# RUN apk add --no-cache docker

# build the app - do I want this to run tests?
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

# run tests
# FROM build-stage AS run-test-stage
# RUN go test -v ./... -count=1

# copy binary to new image
FROM alpine:latest
COPY --from=build-stage /main /main
EXPOSE 23456
ENTRYPOINT ["/main"]
