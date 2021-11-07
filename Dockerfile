# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang as builder

# Copy the local package files to the container's workspace.
COPY . /go/src/instabot
WORKDIR /go/src/instabot
# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get -d -v

RUN go install /go/src/instabot

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/instabot
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy our static executable
COPY --from=builder /go/bin/instabot /go/bin/instabot

# Run the outyet command by default when the container starts.
ENTRYPOINT ["/go/bin/instabot"]