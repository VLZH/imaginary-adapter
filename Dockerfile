FROM golang:alpine

WORKDIR ${GOPATH}/src/github.com/h2non/imaginary
COPY . .
RUN go build

ENV ADAPTER_PORT=9000
ENV ADAPTER_HOST=0.0.0.0
CMD ./imaginary-adapter
