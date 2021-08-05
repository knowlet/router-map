# syntax=docker/dockerfile:experimental

# Set Builder environment
FROM golang:alpine as builder
WORKDIR /workdir
COPY . .
RUN go build -o car src/main.go

# Set deployment environment
FROM alpine as base
WORKDIR /workdir

# Packaging deployment image
FROM base
COPY --from=builder /workdir/car /workdir/car

CMD [ "./car" ]