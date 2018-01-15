FROM golang:1.9.2
LABEL maintainer="TeamHarekaze@harekaze.com"

WORKDIR /go/src/app
COPY src .

RUN cat lib_install.sh | sed 's/go\ get\ -u/go-wrapper\ download/g' | bash
RUN apt-get update && apt-get install -y mysql-client

EXPOSE 5000
CMD ["go", "run", "harekazectf.go"]
