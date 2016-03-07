FROM alpine
ENV MUNKI_REPO_PATH=/data
ENV CADDY_PORT=443

COPY root /

RUN mkdir -p /data/pkgsinfo /data/pkgs /data/catalogs /data/manifests
VOLUME /data
EXPOSE 80 443 2015

CMD ["/ape"]
