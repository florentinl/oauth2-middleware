FROM golang:1.18-alpine3.16
WORKDIR /app

COPY go.mod *.go ./
COPY internal ./internal
RUN CGO_ENABLED=0 go build -o middleware .

FROM alpine
RUN adduser -D middleware
USER middleware
WORKDIR /home/middleware
COPY --from=0 /app/middleware middleware
ENTRYPOINT [ "./middleware" ]
