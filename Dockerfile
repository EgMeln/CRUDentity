FROM golang:latest

WORKDIR /

COPY go.mod go.sum ./

RUN go mod download

COPY main.go /
COPY /internal /internal


RUN go build -o EgMeln/CRUDentity .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache libc6-compat
WORKDIR /root/
COPY --from=0 / .

EXPOSE 8080

CMD ["EgMeln/CRUDentity"]