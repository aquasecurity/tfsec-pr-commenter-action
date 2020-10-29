FROM golang:alpine AS builder

ENV SRCPATH $GOPATH/src/commenter
COPY ./ $SRCPATH
RUN go install commenter


FROM alpine:3.12

RUN apk --no-cache --update add bash git \
    && rm -rf /var/cache/apk/*

SHELL ["/bin/bash", "-eo", "pipefail", "-c"]

COPY --from=builder /go/bin/commenter /

RUN chmod +x /commenter

RUN wget -O - -q "$(wget -q https://api.github.com/repos/tfsec/tfsec/releases/latest -O - | grep -o -E "https://.+?tfsec-linux-amd64")" > tfsec \
    && install tfsec /usr/local/bin/

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]