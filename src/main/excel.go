package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"model/DbModel/StaffModel"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tealeg/xlsx"
)

func indexHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
}

//func uploadHandle(w http.ResponseWriter, r *http.Request) {
func getFile(r *http.Request, w http.ResponseWriter, name string) (string, error) {
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		fmt.Println("ad")
		fmt.Println(token)
		//t, _ := template.ParseFiles("upload.gtpl")
		//t.Execute(w, token)
		return "", nil
	}
	//file, handler, err := r.FormFile("file")
	file, handler, err := r.FormFile(name)
	if err != nil {
		log.Fatal("FormFile: ", err.Error())
		return "", err
	}
	fmt.Fprintf(w, "%v", handler.Header)
	basePath, _ := os.Getwd()
	filePath := basePath + "/upload/" + handler.Filename
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("FormFile: ", err.Error())
		return "", err
	}
	io.Copy(f, file)
	defer func() {
		//关闭文件
		if err := file.Close(); err != nil {
			log.Fatal("Close: ", err.Error())
		}
		if err := f.Close(); err != nil {
			log.Fatal("Close: ", err.Error())
		}
	}()
	return filePath, nil
}

func uploadHandle(w http.ResponseWriter, r *http.Request) {
	path, err := getFile(r, w, "file")
	if err != nil {
		panic(err)
	}
	processFile(path)
}

func main() {
	http.HandleFunc("/upload", uploadHandle)
	http.HandleFunc("/", indexHandle)
	http.ListenAndServe(":9090", nil)
}

func processFile(fileName string) (uint16, error) {
	xlFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		fmt.Printf(err.Error())
	}
	model := StaffModel.NewObj()
	fmt.Println((*model).TabName)
	var cnt uint16 = 0
	for _, sheet := range xlFile.Sheets {
		Rlen := len(sheet.Rows)
		for j := 1; j < Rlen; j++ {
			row := sheet.Rows[j]
			tmp := make(map[string]string)
			len := len(row.Cells)
			if len < 3 {
				continue
			}
			cnt++
			tmp["name"] = "hello"
			tmp["company_id"] = "591591347249199023"
			tmp["id_no"] = row.Cells[2].Value
			tmp["uuid"] = "0"
			tmp["status"] = "0"
			tmp["mobilephone"] = row.Cells[1].Value
			id := model.Insert(tmp)
			fmt.Println(id)
		}
	}
	fmt.Println("done")
	return cnt, nil
}
