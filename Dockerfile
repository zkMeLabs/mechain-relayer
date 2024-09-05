FROM golang:1.22.4-bullseye AS builder
WORKDIR /workspace
COPY . .
RUN make build

FROM golang:1.22.4-bullseye
RUN apt-get update -y && apt-get install ca-certificates -y
COPY --from=builder /workspace/build/greenfield-relayer  /usr/bin/greenfield-relayer

CMD ["/usr/bin/greenfield-relayer"]