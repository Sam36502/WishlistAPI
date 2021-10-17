FROM alpine:latest
RUN mkdir /data
COPY api_linux /data/api_linux
COPY certs /certs
ENV WISHLIST_SSL_CERT /certs/fullchain.pem
ENV WISHLIST_SSL_KEY /certs/privkey.pem
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
CMD ["/data/api_linux"]
