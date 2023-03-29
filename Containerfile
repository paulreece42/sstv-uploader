FROM alpine:3.17

COPY bin/sstv-uploader .

ENTRYPOINT [ "./sstv-uploader" ]
EXPOSE 14230
