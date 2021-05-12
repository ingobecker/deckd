FROM golang:1.16

RUN useradd --create-home go
RUN apt-get update && apt-get install -y libjack-dev
USER go

RUN mkdir /home/go/src
WORKDIR /home/go/src
RUN go install github.com/go-delve/delve/cmd/dlv@latest

LABEL SHELL="podman run --rm -it \
             --name=godev \
             --hostname=godev \
             --userns=keep-id \
             --cap-add=SYS_PTRACE \
             -v .:/home/go/src:Z \
             IMAGE"
