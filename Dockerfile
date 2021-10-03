FROM alpine:latest
RUN mkdir /data
COPY api_linux /data/api_linux
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
CMD ["/data/api_linux"]
