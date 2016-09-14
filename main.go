package main


import(  
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/awsutil"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/aws/session"
    "os"
    "fmt"
    "github.com/kr/fs"
    "net/http"
    "bytes"
    "github.com/solher/arangolite"
    "path/filepath"
    //"time"

)

func main() {  
	db := arangolite.New().
    LoggerOptions(true, true, true).Connect("http://localhost:8529", "personal", "root", "")


    aws_access_key_id := "xxxx"
    aws_secret_access_key := "xxxx"
    token := ""
    creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, token)
    _, err := creds.Get()
    if err != nil {
        fmt.Printf("bad credentials: %s", err)
    }
    
    cfg := aws.NewConfig().WithRegion("us-east-1").WithCredentials(creds)
    svc := s3.New(session.New(), cfg)
    walker := fs.Walk("/Users/vineetdaniel/Desktop")
    //walker := fs.Walk("/Users/vineetdaniel/Documents")
	for walker.Step() {
	    if err := walker.Err(); err != nil {
	        fmt.Fprintln(os.Stderr, err)
	        continue
	    }
	    fmt.Println(walker.Path())
	    path := walker.Path()
	
    file, err := os.Open(path)
    if err != nil {
        fmt.Printf("err opening file: %s", err)
    }

    defer file.Close()

    fileInfo, _ := file.Stat()
    f, err := os.Stat(path)
     if err != nil {
        fmt.Println(err)
     }

    modifiedtime := f.ModTime()
    fmt.Println(fileInfo.Name)
    var size int64 = fileInfo.Size()
    _, name := filepath.Split(path)

    buffer := make([]byte, size)

    // read file content to buffer
    file.Read(buffer)

    fileBytes := bytes.NewReader(buffer)

    fileType := http.DetectContentType(buffer)

    path1 := "/media" + file.Name()
    url := "https://s3.amazonaws.com" + path1
    params := &s3.PutObjectInput{
        Bucket:        aws.String("lpc-v2"),
        Key:           aws.String(path1),
        Body:          fileBytes,
        ContentLength: aws.Int64(size),
        ContentType:   aws.String(fileType),
    }
    if size > 0 {


    resp, err := svc.PutObject(params)
    if err != nil {
        fmt.Printf("bad response: %s", err)
    } else {
    	fmt.Printf("response %s", awsutil.StringValue(resp))
    	q := arangolite.NewQuery("INSERT {path: '%s', orig_path: '%s',size:'%d',name: '%s',fileType: '%s',url: '%s',lastModifiedTime:'%s',uploaded: DATE_NOW()} INTO images",filepath.Clean(path1),path,size,name,fileType,url,modifiedtime)
    	r , err := db.Run(q)
    	fmt.Println(r)
    }
    }
} 
}



