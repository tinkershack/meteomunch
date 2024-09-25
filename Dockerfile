# Build the application from source
FROM golang:1.23.1 AS build-stage

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go vet -v
# RUN go test -v

RUN CGO_ENABLED=0 go build -o /go/bin/app

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /go/bin/app /app

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app"]