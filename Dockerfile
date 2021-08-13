FROM golang:1.16

WORKDIR /go/src/github.com/script-lab/jwt-auth
COPY go.mod .
COPY go.sum .

RUN go get -u -v github.com/cosmtrek/air

COPY . .

CMD ["go", "run", "*.go"]