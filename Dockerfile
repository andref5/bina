FROM ubuntu:20.04

WORKDIR /bina

RUN apt-get update; apt-get install -y clang gcc-multilib vim git curl iproute2

RUN curl -s https://dl.google.com/go/go1.17.linux-amd64.tar.gz | tar -v -C /usr/local -xz
ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV PATH $PATH:/usr/local/go/bin

COPY . .

RUN go build -o http-server http-server.go
CMD ["/bina/http-server"]
