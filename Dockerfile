# Built following https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324

# STEP 1 build executable binary
FROM golang:alpine as builder
# Install SSL ca certificates
RUN apk update && apk add git && apk add ca-certificates
# Create appuser
RUN adduser -D -g '' appuser
COPY . $GOPATH/src/mypackage/myapp/
WORKDIR $GOPATH/src/mypackage/myapp/
#get dependancies
RUN go get -d -v
#build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/acrobits-sipis-exporter


# STEP 2 build a small image
# start from scratch
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable
COPY --from=builder /go/bin/acrobits-sipis-exporter /go/bin/acrobits-sipis-exporter

EXPOSE 8086
USER appuser
ENTRYPOINT ["/go/bin/acrobits-sipis-exporter", "-listen-address", ":8086"]
