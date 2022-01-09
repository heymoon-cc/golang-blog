FROM golang:1.17 as builder
WORKDIR /code
ADD ./go.* /code/
RUN go mod download
ADD ./*.go /code/
ADD ./controller /code/controller
ADD ./model /code/model
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o blog

FROM busybox as base
WORKDIR /etc/blog
COPY --from=builder /code/blog /
ADD ./ui /etc/blog/ui
ADD ./static /etc/blog/static
ADD ./config.json /etc/blog/
CMD ["/blog"]
