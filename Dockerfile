FROM golang:1.13 as go-builder

RUN mkdir -p /go/src/github.com/nitronick600/nhvolume
WORKDIR /go/src/github.com/nitronick600/nhvolume

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/nhvolume

FROM scratch

WORKDIR /
COPY --from=go-builder /go/bin/nhvolume /go/bin/nhvolume
COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/go/bin/nhvolume"]
