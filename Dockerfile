FROM jheise/ubuntu-golang

RUN mkdir -p /go/src/dockerview
ADD *.go /go/src/dockerview
RUN go get dockerview
RUN go install dockerview
ADD static /srv/static
ADD templates /srv/templates
ENV ADDRESS 0.0.0.0
ENV PORT 9999
EXPOSE ${PORT}
WORKDIR /srv
CMD /go/bin/dockerview -address ${ADDRESS} -port ${PORT}
