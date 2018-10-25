FROM golang

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./... && go install -v ./...

FROM debian

RUN apt-get update && apt-get install -y ca-certificates
COPY --from=0 /go/bin/app /app
CMD ["/app"]
