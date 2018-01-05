FROM golang:1.9.2
MAINTAINER HayatoDoi <nono@nononono.net>

WORKDIR /go/src/app
COPY src .

RUN cat lib_install.sh | sed 's/go\ get\ -u/go-wrapper\ download/g' | bash

CMD ["go", "run", "harekazectf.go"]
