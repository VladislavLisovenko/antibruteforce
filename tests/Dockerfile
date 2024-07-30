FROM golang:1.22-alpine3.19 as build

ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY ./tests/go.mod .
COPY ./tests/go.sum .
RUN ls -la ./

RUN go mod download

COPY ./tests .

CMD ["go", "test"]