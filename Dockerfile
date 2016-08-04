FROM golang:1.6

ADD . /go/src/github.com/briankassouf/learn
WORKDIR /go/src/github.com/briankassouf/learn
RUN go get github.com/mattn/goveralls
RUN go get github.com/tools/godep
RUN godep restore
RUN godep go install -v .

CMD ["learn"]
