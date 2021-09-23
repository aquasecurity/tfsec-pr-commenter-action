FROM alpine:3.12

RUN apk --no-cache --update add bash git

RUN wget https://gitreleases.dev/gh/aquasecurity/tfsec/latest/tfsec-linux-arm64 > tfsec \
    && install tfsec /usr/local/bin/

RUN wget https://gitreleases.dev/gh/aquasecurity/tfsec-pr-commenter-action/latest/commenter-linux-amd64 > commenter \
    && install commenter /usr/local/bin/

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
