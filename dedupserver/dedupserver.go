
package main

import "net/http"
import "fmt"
import "os"
import "io"
import "path/filepath"
import "strconv"
import "sync"
import "crypto/sha256"

//list of known files
//key = original filename; value = its integer index
var filemap map[string] int;

//number of files we know about
var numfiles int;
var fileLock sync.Mutex;

//stored file hashes
var fileHashes map[string] int;

func main(){
    filemap = make(map[string]int)
    fileHashes = make(map[string]int)
	
    
    //create a directory for our data. 
    //mode = 0777 = rwxrwxrwx
    os.RemoveAll("dedupdata")
    os.Mkdir("dedupdata",os.ModeDir | 0777);
    
    http.HandleFunc("/", indexPage )
    http.HandleFunc("/upload", uploadData )
    http.HandleFunc("/download", downloadData )
    
    fmt.Println("Listening on port 8989");
    
    //this is automatically multithreaded: Each client
    //gets a different goroutine
    err := http.ListenAndServe(":8989", nil)
    if err != nil {
        fmt.Println("Cannot listen");
    }
}

func indexPage(w http.ResponseWriter, r *http.Request){
    fmt.Println(r.Method,r.URL)
    fmt.Fprintln(w,"<!DOCTYPE html>")
    fmt.Fprintln(w,"<HTML><head><meta charset=utf8></head>")
    fmt.Fprintln(w,"<body>")
    fmt.Fprintln(w,"<form method='post' action='upload' enctype='multipart/form-data'>")
    fmt.Fprintln(w,"<input type=file name=uploadfile>")
    fmt.Fprintln(w,"<input type=submit value='Upload it!'>")
    fmt.Fprintln(w,"</form>")
    fmt.Fprintln(w,"File list:<br>");
	fileLock.Lock()
    for filename, filenumber := range filemap {
        fmt.Fprintf(w,"<a href='/download?number=%d'>%s</a><br>\n" , 
            filenumber, filename)
    }
	fileLock.Unlock()
    fmt.Fprintln(w,"</body>");
    fmt.Fprintln(w,"</html>");
}

func uploadData( w http.ResponseWriter, r *http.Request){
    f,fh,err := r.FormFile("uploadfile")
    if err != nil {
        fmt.Println("Error when uploading",err)
        return;
    }
    filename := fh.Filename
	
	hash := sha256.New()
	io.Copy(hash, f)
	var fileHashSig []byte
	fileHashSig = hash.Sum(fileHashSig)
	f.Seek(0,0) //reset the file because for some reason it doesn't already
	
	//fmt.Println(string(fileHashSig))
	fileLock.Lock()
	fileNum,fileFound := fileHashes[string(fileHashSig)]
	fileLock.Unlock()
	if fileFound{
		fileLock.Lock()
		filemap[filename] = fileNum
		fileLock.Unlock()
		fmt.Println("Hash found in array. Using existing filenumber.")
	}else{
		fileLock.Lock()
		filenumber := numfiles
		numfiles += 1
		fileHashes[string(fileHashSig)] = filenumber
		filemap[filename] = filenumber
		fileLock.Unlock()
		fmt.Println("Hash not found in array. Adding...")
		
		fp,err := os.Create(filepath.Join("dedupdata",strconv.Itoa(filenumber)))
		if err != nil {
			fmt.Println("Cannot create file. Why?")
			return
		}
		
		defer fp.Close()
		_,err = io.Copy(fp,f)
		if err != nil {
			fmt.Println("Error on copy")
		}
		
	}
    indexPage(w,r);
}

func downloadData( w http.ResponseWriter, r *http.Request){
    numberarray := r.URL.Query()["number"]
    if len(numberarray) == 0 {
        fmt.Println("Bad request")
        return
    }
    filenumber := numberarray[0]
    fp,err := os.Open(filepath.Join("dedupdata",filenumber))
    if err != nil {
        fmt.Println("Cannot open file",filenumber)
        w.WriteHeader(404)
        return
    }
    defer fp.Close()
    _,err = io.Copy(w,fp)
    if err != nil {
        fmt.Println("Error on copy")
    }
}
