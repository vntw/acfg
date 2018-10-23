FROM alpine:3.8

COPY build/acfg-linux-amd64 /

ENV ACFG_PORT=1337
ENV ACFG_ACSERVER_DIR=/acserver
ENV ACFG_SERVER_LOGS_DIR=/logs
ENV ACFG_SERVER_CFGS_DIR=/cfgs

VOLUME ["/acserver", "/logs", "/cfgs"]

EXPOSE 1337

CMD /acfg-linux-amd64
