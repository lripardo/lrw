FROM golang:1.17 as builder
LABEL stage=intermediate
ARG VERSION
WORKDIR /go/src/github.com/lripardo/lrw
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.VERSION=${VERSION}" -o lrw

FROM scratch
COPY --from=builder /go/src/github.com/lripardo/lrw/lrw .
EXPOSE 8080
ENTRYPOINT ["./lrw"]
