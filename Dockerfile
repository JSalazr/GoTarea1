FROM golang:1.8

RUN git clone https://github.com/Vacster/Languages_go.git

EXPOSE 8080

CMD cd /go/Languages_go/ && go run server.go
