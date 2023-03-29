FROM alpine:3.17

WORKDIR /usr/src/app

COPY bin/sstv-uploader .

ENTRYPOINT [ "./sstv-uploader" ]
EXPOSE 14230
