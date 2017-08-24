FROM golang:1.8 as builder
WORKDIR /go/src/web
COPY web.go .
COPY vendor ./vendor
RUN go build -v -tags netgo .

FROM scratch
USER 1001
EXPOSE 8080
COPY --from=builder /go/src/web/web /
CMD ["/web"]
