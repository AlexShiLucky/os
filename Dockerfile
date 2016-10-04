FROM alpine:3.2
WORKDIR /
# template dirs
ADD discovery/templates /discovery/templates
ADD config/templates /config/templates
ADD monitor/templates /monitor/templates
ADD event/templates /event/templates
ADD router/templates /router/templates
ADD trace/templates /trace/templates
# binary
ADD os /os
ENTRYPOINT [ "/os" ]
