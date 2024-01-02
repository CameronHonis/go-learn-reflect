FROM golang:bookworm

WORKDIR main

COPY . .

CMD ["./init"]