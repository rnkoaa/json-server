FROM golang:1.14.0 as gobuilder

WORKDIR /app

COPY go.mod .
COPY go.sum .

ENV GO111MODULE=on
# RUN CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o json-server

RUN go mod download

# FROM scratch
# COPY --from=gobuilder /go-app/json-server .
# ENTRYPOINT ["/json-server"]
