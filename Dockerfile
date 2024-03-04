# Builder
FROM golang:1.20.7-alpine3.17 as builder

RUN apk update && apk upgrade && \
    apk --update add git make bash build-base

WORKDIR /app

COPY . .

RUN make deps
RUN make build

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app 

WORKDIR /app 

EXPOSE 8080

COPY --from=builder /app/main /app/
COPY --from=builder /app/.env /app/

CMD /app/main