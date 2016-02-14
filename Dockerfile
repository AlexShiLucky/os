FROM alpine:3.2
WORKDIR /
ADD platform /platform
ENTRYPOINT [ "/platform" ]
