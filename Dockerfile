FROM golang:latest
WORKDIR /app
ADD . .
RUN go mod download
RUN make build
CMD ./out/bin/ms-arch-example