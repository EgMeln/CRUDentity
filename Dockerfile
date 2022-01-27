FROM golang:latest

WORKDIR /

COPY go.mod go.sum ./

RUN go mod download

COPY main.go /
COPY /internal /internal


RUN go build -o EgMeln/CRUDentity .

EXPOSE 8080

CMD ["EgMeln/CRUDentity"]