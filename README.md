
A pointlessly over-complicated way to upload SSTV images from amateur radio that no one should use

(Basically my "hello world" golang project :) )

Nevertheless, if you do decide to run it, you'll need to set some env vars:


    DB_USER=sstv DB_PASS=sstv DB_NAME=sstv S3_REGION=default S3_BUCKET=sstv S3_ENDPOINT=objects.servercloud.com BEARER_TOKEN=foobarbaz go run main.go

Setup your s3 credentials in ~/.aws/credentials

And post images from MMSSTV/QSSTV using something like Facebook's watchman or Linux's inotifywait or something:

    curl -XPOST -H 'Bearer: foobarbaz' -H 'Content-Type: multipart/form-data' --form 'file=@Hist6.bmp' localhost:14230/sstv/



Thanks to https://www.section.io/engineering-education/build-a-rest-api-application-using-golang-and-postgresql-database/ and https://stackoverflow.com/questions/40684307/how-can-i-receive-an-uploaded-file-using-a-golang-net-http-server for getting me started, hopefully I didn't forget too many other sources :-/
