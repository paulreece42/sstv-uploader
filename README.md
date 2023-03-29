
A pointlessly over-complicated way to upload SSTV images from amateur radio that no one should use

(Basically my "hello world" golang project :) )

Nevertheless, if you do decide to run it, see BUILDING


And post images from MMSSTV/QSSTV using something like Facebook's watchman or Linux's inotifywait or something:

    curl -XPOST -H 'Bearer: foobarbaz' -H 'Content-Type: multipart/form-data' --form 'file=@Hist6.bmp' localhost:14230/sstv/

