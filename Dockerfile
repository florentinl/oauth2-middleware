FROM golang:1.17-alpine3.15
WORKDIR /app

COPY go.mod *.go ./
RUN CGO_ENABLED=0 go build -o forward .

FROM alpine
RUN adduser -D forward
USER forward
WORKDIR /home/forward
COPY --from=0 /app/forward forward
ENTRYPOINT [ "./forward" ]
