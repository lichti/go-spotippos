FROM golang:1.8.3 as builder
WORKDIR /go/src/github.com/lichti/go-spotippos/
COPY ./    .
WORKDIR /go/
RUN go get -v -d github.com/lichti/go-spotippos
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o spotippos .
RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo -v github.com/lichti/go-spotippos

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/bin/go-spotippos .
RUN chmod +x ./go-spotippos
CMD ["./go-spotippos"]
EXPOSE 8000
