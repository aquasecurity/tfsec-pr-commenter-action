FROM alpine:3.12

RUN apk --no-cache --update add bash git \
    && rm -rf /var/cache/apk/*

RUN wget -O - -q "$(wget -q https://api.github.com/repos/tfsec/tfsec/releases/latest -O - | grep -o -E "https://.+?tfsec-linux-amd64")" > tfsec \
    && install tfsec /usr/local/bin/

RUN wget -O - -q "$(wget -q https://github.com/tfsec/tfsec-pr-commenter-action/releases/latest -O - | grep -o -E "https://.+?commenter-linux-amd64")" > commenter \
    && install commenter /usr/local/bin/

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
