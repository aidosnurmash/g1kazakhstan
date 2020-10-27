package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	_ "image/jpeg"
	_ "image/png"
	"imageSavingProject/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)
var port string
var server string
var dbPort string
var dbServer string
var filePath string
var database *models.Database

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/images/{id:[0-9]+}", viewImageHandler)
	router.HandleFunc("/images/{id:[0-9]+}/part/{part_num:[1-4]}", viewImagePartHandler)
	router.HandleFunc("/images", saveImageHandler).Methods("POST")
	router.HandleFunc("/", mainPageHandler)
	http.Handle("/",router)

	flag.StringVar(&port, "port", "8000", "specify port to use.  defaults to 8000")
	flag.StringVar(&server, "server", "127.0.0.1", "specify server to use.  defaults to 127.0.0.1")
	flag.StringVar(&dbPort, "db_port", "8000", "specify database port to use.  defaults to 8000")
	flag.StringVar(&dbServer, "db_server", "127.0.0.1", "specify dbServer to use.  defaults to 127.0.0.1")
	flag.StringVar(&filePath, "file_path", "./images", "specify filePath to use.  defaults to ./images")
	flag.Parse()

	fmt.Println(port)
	fmt.Println(server)
	fmt.Println(dbPort)
	fmt.Println(dbServer)
	fmt.Println(filePath)
	fmt.Println(os.Args)
	filePath = "/Users/aidos/go/src/imageSavingProject/images/"

	database = &models.Database{}
	database.Init()
	log.Fatal(http.ListenAndServe(server+":"+port, nil))
	defer database.Db.Close()
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w,"<!DOCTYPE html>\n<html lang=\"en\">\n  <head>\n    <meta charset=\"UTF-8\" />\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\n    <meta http-equiv=\"X-UA-Compatible\" content=\"ie=edge\" />\n    <title>Document</title>\n  </head>\n  <body>\n    <form\n      enctype=\"multipart/form-data\"\n      action=\"/images\"\n      method=\"post\"\n    >\n      <input type=\"file\" name=\"image\" />\n      <input type=\"submit\" value=\"upload\" />\n    </form>\n  </body>\n</html>")
}

func saveImageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "I'm in saveImageHandler")

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)




	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile, err := ioutil.TempFile(filePath, "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()
	tempFile.Write(fileBytes)
	/*imageData, imageType, err := image.Decode(file)
	if err != nil {
		fmt.Fprintf(w, "Error  opening image %+v\n", err.Error())
		return
	}
	fmt.Println("okasdfasd")
	*/
	insertedPictureId, _ := database.InsertPicture(tempFile.Name(), handler.Filename)


	fmt.Fprintf(w, "Successfully Uploaded File %+v %+v\n", tempFile.Name(), insertedPictureId)
}

func saveFileByte(fileBytes []byte) string {
	tempFile, err := ioutil.TempFile(filePath, "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()
	tempFile.Write(fileBytes)
	return tempFile.Name()
}

func viewImageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("OK")
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	fmt.Println(id)
	picture, err := database.GetPictureById(id)
	if err != nil {
		//fmt.Fprintln(w, "Image not exist")
		return
	}
	fmt.Println(picture)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err) // perhaps handle this nicer
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote("download.png"))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, picture.Path)
}
func viewImagePartHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("OK")
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	part_num, _ := strconv.Atoi(vars["id"])
	fmt.Println(id)
	picture, err := database.GetPictureById(id)
	if err != nil {
		fmt.Println(w, "picture not exist")
		return
	}
	part, err := database.GetPartById(picture.Id, part_num)
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote("download.png"))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, part.Path)
}
