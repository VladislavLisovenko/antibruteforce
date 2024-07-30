FROM golang:1.22-alpine3.19 as build

ENV BIN_FILE /opt/anti-bruteforce/anti-bruteforce
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/anti-bruteforce/*

FROM alpine:3.19

ENV BIN_FILE "/opt/anti-bruteforce/anti-bruteforce"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

CMD ${BIN_FILE}