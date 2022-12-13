FROM golang:alpine as buildenv
ADD .   /app
WORKDIR /app
RUN mkdir -p /app/bin
RUN apk --no-cache --update add gcc musl-dev
RUN go build -mod=vendor -o bin/cart github.com/cubny/cart/cmd/...

#--------

FROM alpine:latest
COPY --from=buildenv /app/bin/* /app/bin/
EXPOSE 8080 8081
CMD ["/app/bin/cart"]
