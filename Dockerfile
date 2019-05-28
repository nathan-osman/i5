FROM golang:latest

# Disable cgo, to ensure a standalone binary
ENV CGO_ENABLED=0

ADD . /src

WORKDIR /src

RUN go build


FROM scratch

COPY --from=0 /src/i5 /usr/local/bin/

# Use cURL's certificate bundle
ADD https://curl.haxx.se/ca/cacert.pem /etc/ssl/certs/

ENTRYPOINT ["/usr/local/bin/i5"]
