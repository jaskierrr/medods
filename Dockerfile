# syntax=docker/dockerfile:1

FROM golang:latest as builder

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .


ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux go build -v -o /medods ./cmd/app/main.go


FROM alpine:latest

COPY --from=builder medods /bin/main
COPY .env .


EXPOSE 8080

ENTRYPOINT ["/bin/main"]
