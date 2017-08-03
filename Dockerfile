FROM golang:1.8.3 as builder
WORKDIR /go/src/github.com/lichti/spotippos/
COPY ./    .
WORKDIR /go/
RUN go get -v -d github.com/lichti/spotippos
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o spotippos .
RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo -v github.com/lichti/spotippos

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/bin/spotippos .
RUN chmod +x ./spotippos
CMD ["./spotippos"]