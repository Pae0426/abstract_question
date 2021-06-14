package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	//"shares/package/pdf_to_image"
)

//指定ディレクトリ下のファイル数をカウントする
func countFiles() int {
	files, _ := ioutil.ReadDir("./static/pdf/1/")
	var count int
	for _, f := range files {
		if f.Name() == ".DS_Store" {
			continue
		}
		count++
	}
	return count
}

func templateHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]int{
		"pages": countFiles(),
	}
	t, err := template.ParseFiles(
		"/go/src/app/views/home.html",
		"/go/src/app/views/header.html",
		"/go/src/app/views/footer.html",
		//"views/home.html",
		//"views/header.html",
		//"views/footer.html",
	)
	if err != nil {
		log.Fatalln("テンプレートファイルを読み込めません:", err.Error())
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatalln("エラー!:", err.Error())
	}
}

/*func convertHandler(w http.ResponseWriter, r *http.Request) {
	pdf_to_image.Pdf_to_image()
}*/

func main() {
	log.Println("Webサーバーを開始します...")
	server := http.Server{
		Addr: ":8080",
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/home", templateHandler)
	//http.HandleFunc("/pdf_to_image", convertHandler)
	server.ListenAndServe()
}
