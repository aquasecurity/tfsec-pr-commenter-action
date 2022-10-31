FROM public.ecr.aws/docker/library/alpine:3.12

RUN apk --no-cache --update add bash git \
    jq curl \
    && rm -rf /var/cache/apk/*

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
