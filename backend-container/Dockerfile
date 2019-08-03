FROM golang:alpine
ADD ./src /app
WORKDIR /app

RUN apk add --update tzdata \
    bash wget curl git gcc libc-dev;

RUN export GO111MODULE=on

RUN go get github.com/oxequa/realize

EXPOSE 80

CMD realize start --run serve