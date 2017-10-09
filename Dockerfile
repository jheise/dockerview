FROM jheise/ubuntu-golang

RUN mkdir -p /go/src/dockerview
ADD *.go /go/src/dockerview
RUN go get dockerview
RUN go install dockerview
ADD static /srv/static
ADD templates /srv/templates
EXPOSE 9999
WORKDIR /srv
CMD /go/bin/dockerview
