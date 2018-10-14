FROM golang:1.10
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep


RUN mkdir -p /go/src/github.com/jack-slater/go-login
COPY . /go/src/github.com/jack-slater/go-login
WORKDIR /go/src/github.com/jack-slater/go-login/app

RUN wget https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh
RUN chmod +x wait-for-it.sh

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

RUN go build -o main .