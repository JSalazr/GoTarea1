FROM golang:1.8

RUN git clone https://github.com/JSalazr/GoTarea1.git
RUN go get googlemaps.github.io/maps
RUN go get golang.org/x/net/context
RUN go get golang.org/x/image/bmp
RUN go get github.com/kr/pretty

EXPOSE 8080

CMD cd /go/GoTarea1 && go run Server.go
