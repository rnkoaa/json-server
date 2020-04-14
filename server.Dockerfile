FROM golang:1.14.0 as gobuilder

WORKDIR /app

COPY go.mod .
COPY go.sum .

ENV GO111MODULE=on

RUN go mod download

# copy the contents of the image and build it.
FROM gobuilder as builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o json-server ./cmd/server/...

# now create a minimal image that will build the application.
FROM alpine:3.7
EXPOSE 8080
WORKDIR /app/
COPY --from=builder /app/json-server .
COPY --from=builder /app/cmd/server/db.json .

ENTRYPOINT ["./json-server"]
