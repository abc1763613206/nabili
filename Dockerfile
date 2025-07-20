FROM golang:alpine as builder

RUN apk add --no-cache make git
WORKDIR /nabili-src
COPY . /nabili-src
RUN go mod download && \
    make docker && \
    mv ./bin/nabili-docker /nabili

FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /nabili /
ENTRYPOINT ["/nabili"]