FROM golang:1.9.2
LABEL maintainer="TeamHarekaze@harekaze.com"

RUN mkdir -p $GOPATH/src/github.com/HayatoDoi/HarekazeCTF-Competition
WORKDIR $GOPATH/src/github.com/HayatoDoi/HarekazeCTF-Competition
COPY . .

RUN cat lib_install.sh | sed 's/go\ get\ -u/go-wrapper\ download/g' | bash
RUN apt-get update && apt-get install -y mysql-client

EXPOSE 5000
CMD ["go", "run", "main.go"]
