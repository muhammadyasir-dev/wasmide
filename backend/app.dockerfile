FROM golang:1.23.4 AS builder

WORKDIR /backend

COPY go.mod go.sum

COPY . .

#EXPOSE 8080

#CMD["./main"]

RUN go build -o main .

FROM golang:1.23.4

WORKDIR /backend


COPY --from=bulder /backend/main .

EXPOSE 8080

CMD["./main"]
