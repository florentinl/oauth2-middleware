FROM golang:1.17-alpine3.15
RUN mkdir /app
COPY go.mod kubernetes-oauth2-forward.go /app/
WORKDIR /app
RUN go build -o main .
ENTRYPOINT [ "/app/main" ]