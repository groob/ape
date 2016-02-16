FROM alpine
COPY ./bin/ape-linux-amd64 /ape

ENV MUNKI_REPO_PATH=/data
ENV APE_HTTP_LISTEN_PORT=80
VOLUME /data
EXPOSE 80

CMD ["/ape"]
