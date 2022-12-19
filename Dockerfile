FROM golang:1.19.4-alpine as builder
RUN apk --update add ca-certificates
RUN cd ..
RUN mkdir httpqueue
WORKDIR httpqueue
COPY . ./
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o httpqueue ./cmd/main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/httpqueue/httpqueue .
CMD ["./httpqueue"]
