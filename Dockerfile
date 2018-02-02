FROM golang:1.9.3-alpine

ENV GOMAXPROCS 1
ENV GOPATH /root/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ADD ./ $GOPATH/src/github.com/oleggator/esports-backend

CMD go install -ldflags '-s' github.com/oleggator/esports-backend && tp-db
