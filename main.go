package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "bytes"
    "io"
    "time"
    "os"
    "strconv"
    "crypto/md5"
    "encoding/hex"

    "image"
    "image/png"
    "golang.org/x/image/bmp"

    "github.com/gorilla/mux"
    "github.com/google/uuid"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    _ "github.com/lib/pq"
)


var (
    DB_USER string = os.Getenv("DB_USER")
    DB_PASSWORD string = os.Getenv("DB_PASS")
    DB_NAME string = os.Getenv("DB_NAME")
    DB_HOST string = os.Getenv("DB_HOST")

    S3_REGION string = os.Getenv("S3_REGION")
    S3_BUCKET string = os.Getenv("S3_BUCKET")
    S3_ENDPOINT string = os.Getenv("S3_ENDPOINT")
    BEARER_TOKEN string = os.Getenv("BEARER_TOKEN")
//    S3_TOKEN string = os.Getenv("S3_TOKEN") // set these in ~/.aws/credentials
//    S3_PRIVKEY string = os.Getenv("S3_PRIVKEY")

)


func printMessage(message string) {
    fmt.Println("")
    fmt.Println(message)
    fmt.Println("")
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}


func setupDB() *sql.DB {
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s", DB_USER, DB_PASSWORD, DB_NAME, DB_HOST)
    db, err := sql.Open("postgres", dbinfo)

    checkErr(err)

    return db
}

type SSTV struct {
    SSTV_ID   string `json:"sstvid"`
    UploadTime string `json:"uploadtime"`
    Link string `json:"link"`
}

type JsonResponse struct {
    Type    string `json:"type"`
    Data    []SSTV `json:"data"`
    Message string `json:"message"`
}

var db = setupDB()

func main() {

    // Init the mux router
    router := mux.NewRouter()

// Route handles & endpoints

    router.HandleFunc("/sstv/", GetSSTV).Methods("GET")
    router.HandleFunc("/sstv/", CreateSSTV).Methods("POST")
    router.HandleFunc("/sstv/{page}", GetSSTVPage).Methods("GET")
//    router.HandleFunc("/sstv/{sstvid}", DeleteSSTV).Methods("DELETE")


    fmt.Println("Server at 14230")
    log.Fatal(http.ListenAndServe(":14230", router))
}


func GetSSTV(w http.ResponseWriter, r *http.Request) {
    bearer := r.Header.Get("Bearer")
    if ( bearer != BEARER_TOKEN) {
        response := JsonResponse{Type: "error", Message: "Must set bearer token header"}
        json.NewEncoder(w).Encode(response)
        return
    }

    

    printMessage("Getting sstv...")
    rows, err := db.Query("SELECT sstvid, uploaded_at FROM sstv order by uploaded_at desc")
    checkErr(err)
    defer rows.Close()

    var SSTVs []SSTV

    s, err := session.NewSession(&aws.Config{Region: aws.String(S3_REGION), Endpoint: aws.String(S3_ENDPOINT)})
    checkErr(err)


    mys3 := s3.New(s)



    for rows.Next() {
        var sstv_id string
        var uploadtime string

        err = rows.Scan(&sstv_id, &uploadtime)

        checkErr(err)

        req, _ := mys3.GetObjectRequest(&s3.GetObjectInput{
            Bucket: aws.String(S3_BUCKET),
            Key:    aws.String(sstv_id + ".png"),
        })
        link, err := req.Presign(120 * time.Minute)

        checkErr(err)

        SSTVs = append(SSTVs, SSTV{SSTV_ID: sstv_id, UploadTime: uploadtime, Link: link})
    }

    var response = JsonResponse{Type: "success", Data: SSTVs}
    json.NewEncoder(w).Encode(response)
}

func GetSSTVPage(w http.ResponseWriter, r *http.Request) {
    bearer := r.Header.Get("Bearer")
    if ( bearer != BEARER_TOKEN) {
        response := JsonResponse{Type: "error", Message: "Must set bearer token header"}
        json.NewEncoder(w).Encode(response)
        return
    }
    params := mux.Vars(r)
    page, err := strconv.Atoi(params["page"])


    printMessage("Getting sstv...")
    rows, err := db.Query("SELECT sstvid, uploaded_at FROM sstv order by uploaded_at desc limit 10 offset $1", page * 10)
    checkErr(err)
    defer rows.Close()

    var SSTVs []SSTV

    s, err := session.NewSession(&aws.Config{Region: aws.String(S3_REGION), Endpoint: aws.String(S3_ENDPOINT)})
    checkErr(err)


    mys3 := s3.New(s)



    for rows.Next() {
        var sstv_id string
        var uploadtime string

        err = rows.Scan(&sstv_id, &uploadtime)

        checkErr(err)

        req, _ := mys3.GetObjectRequest(&s3.GetObjectInput{
            Bucket: aws.String(S3_BUCKET),
            Key:    aws.String(sstv_id + ".png"),
        })
        link, err := req.Presign(120 * time.Minute)

        checkErr(err)

        SSTVs = append(SSTVs, SSTV{SSTV_ID: sstv_id, UploadTime: uploadtime, Link: link})
    }

    var response = JsonResponse{Type: "success", Data: SSTVs}
    json.NewEncoder(w).Encode(response)
}

func CreateSSTV(w http.ResponseWriter, r *http.Request) {
    bearer := r.Header.Get("Bearer")
    if ( bearer != BEARER_TOKEN) {
        response := JsonResponse{Type: "error", Message: "Must set bearer token header"}
        json.NewEncoder(w).Encode(response)
        return
    }
    r.ParseMultipartForm(32 << 20) // limit your max input length!
    var buf bytes.Buffer
    // in your case file would be fileupload
    file, header, err := r.FormFile("file")
    checkErr(err)

    name := header.Filename
    fmt.Printf("File name %s\n", name[0])

    io.Copy(&buf, file)

    myuuid := uuid.New().String()
    defer file.Close()

    s, err := session.NewSession(&aws.Config{Region: aws.String(S3_REGION), Endpoint: aws.String(S3_ENDPOINT)})
    checkErr(err)


    hash := md5.Sum(buf.Bytes())
    mysum := hex.EncodeToString(hash[:])

    fmt.Printf("md5sum: %s\n", mysum)

    var cnt int

    rows, err := db.Query("SELECT count(*) as cnt FROM sstv where md5sum = $1", mysum)
    checkErr(err)
    rows.Next()
    rows.Scan(&cnt)
    if (cnt > 0) {
        response := JsonResponse{Type: "error", Message: "Content already exists"}
        json.NewEncoder(w).Encode(response)
        return
    }
    fmt.Printf("return %d", cnt)
    mys3 := s3.New(s)
    defer rows.Close()

    var size int64 = int64(buf.Len())

    _, err = mys3.PutObject(&s3.PutObjectInput{
        Bucket:               aws.String(S3_BUCKET),
        Key:                  aws.String(myuuid),
        ACL:                  aws.String("private"),
        Body:                 bytes.NewReader(buf.Bytes()),
        ContentLength:        aws.Int64(size),
        ContentType:          aws.String(http.DetectContentType(buf.Bytes())),
    })
    checkErr(err)


    if ( http.DetectContentType(buf.Bytes()) == "image/bmp" ) {
        var img image.Image
        img, err = bmp.Decode(bytes.NewReader(buf.Bytes()))
        buf.Reset()
        png.Encode(&buf, img)

        size  = int64(buf.Len())

        _, err = mys3.PutObject(&s3.PutObjectInput{
            Bucket:               aws.String(S3_BUCKET),
            Key:                  aws.String(myuuid + ".png"),
            ACL:                  aws.String("private"),
            Body:                 bytes.NewReader(buf.Bytes()),
            ContentLength:        aws.Int64(size),
            ContentType:          aws.String(http.DetectContentType(buf.Bytes())),
    })
    checkErr(err)


    }

    var response = JsonResponse{}

    response = JsonResponse{Type: "success", Message: myuuid}

    

    printMessage("Inserting SSTV into DB")

    fmt.Println("Inserting new SSTV with ID: " + myuuid)

    var lastInsertID string
    err = db.QueryRow("insert into sstv (sstvid, md5sum) values ($1, $2) returning sstvid", myuuid, mysum).Scan(&lastInsertID)

    // check errors
    checkErr(err)

    response = JsonResponse{Type: "success", Message: "Successfully posted: " + lastInsertID}
    json.NewEncoder(w).Encode(response)
}
