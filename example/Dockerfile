FROM golang:alpine

WORKDIR /app
COPY . /app
RUN apk add git --no-cache
RUN go get ./... && go build .

FROM alpine

COPY --from=0 /app/app /app
CMD ["/app"]
