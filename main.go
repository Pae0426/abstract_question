package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Sticky struct {
	Id       int    `json:"id,omitempty"`
	Page     int    `json:"page,omitempty"`
	Color    string `json:"color,omitempty"`
	Shape    string `json:"shape,omitempty"`
	Locate_x int    `json:"location_x,omitempty"`
	Locate_y int    `json:"location_y,omitempty"`
	Text     string `json:"text,omitempty"`
	Empathy  int    `json:"empathy,omitempty"`
}

var Db *sql.DB

func init() {
	var err error
	err = godotenv.Load(".env")
	if err != nil {
		log.Println("エラー:", err)
	}

	DB_NAME := os.Getenv("MYSQL_DATABASE")
	DB_USER := os.Getenv("MYSQL_USER")
	DB_PASS := os.Getenv("MYSQL_PASSWORD")
	DB_PROTOCOL := "tcp(127.0.0.1:3306)"
	DB_CONNECT_INFO := DB_USER + ":" + DB_PASS + "@" + DB_PROTOCOL + "/" + DB_NAME
	Db, err = sql.Open("mysql", DB_CONNECT_INFO)
	if err != nil {
		log.Println("エラー:", err)
	}
}

//ログをファイルに出力
func loggingSettings(filename string) {
	logfile, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetOutput(multiLogFile)
}

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

//付箋の情報をDBから取得しjson形式で表示
func getStickiesInfo(w http.ResponseWriter, r *http.Request) {
	rows, e := Db.Query("select * from lecture1")
	if e != nil {
		log.Println("エラー:", e.Error())
	}

	var stickies []Sticky

	for rows.Next() {
		sticky := Sticky{}
		if er := rows.Scan(
			&sticky.Id,
			&sticky.Page,
			&sticky.Color,
			&sticky.Shape,
			&sticky.Locate_x,
			&sticky.Locate_y,
			&sticky.Text,
			&sticky.Empathy,
		); er != nil {
			log.Println(er)
		}
		stickies = append(stickies, sticky)
	}

	defer rows.Close()

	result, err := json.Marshal(stickies)
	if err != nil {
		log.Println("エラー:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func loadStickyId(w http.ResponseWriter, r *http.Request) {
	row, e := Db.Query("select max(id) from lecture1")
	if e != nil {
		log.Println("エラー:", e.Error())
	}

	defer row.Close()

	var id int
	for row.Next() {
		if er := row.Scan(&id); er != nil {
			log.Println(er)
		}
	}

	result, err := json.Marshal(id)
	if err != nil {
		log.Println("エラー:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func createSticky(w http.ResponseWriter, r *http.Request) {
	var sticky Sticky
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	if err := json.Unmarshal(body[:len], &sticky); err != nil {
		log.Fatalln("エラー")
	}

	sql, err := Db.Prepare("insert into lecture1(page, color, shape, location_x, location_y, text, empathy) values(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("エラー:", err)
	}
	sql.Exec(sticky.Page, sticky.Color, sticky.Shape, sticky.Locate_x, sticky.Locate_y, sticky.Text, sticky.Empathy)

	res, err := json.Marshal("{200, \"ok\"}")
	if err != nil {
		log.Println("エラー:", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func updateSticky(w http.ResponseWriter, r *http.Request) {
	var sticky Sticky
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	if err := json.Unmarshal(body[:len], &sticky); err != nil {
		log.Println("エラー:", err)
	}

	sql, err := Db.Prepare("update lecture1 set location_x=?, location_y=? where id=?")
	if err != nil {
		log.Println("エラー:", err)
	}
	sql.Exec(sticky.Locate_x, sticky.Locate_y, sticky.Id)

	res, err := json.Marshal("{200, \"ok\"}")
	if err != nil {
		log.Println("エラー:", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func templateHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]int{
		"pages": countFiles(),
	}
	t, err := template.ParseFiles(
		"views/home.html",
		"views/header.html",
		"views/footer.html",
	)
	if err != nil {
		log.Fatalln("テンプレートファイルを読み込めません:", err.Error())
	}
	if err := t.Execute(w, data); err != nil {
		log.Println("エラー:", err.Error())
	}
}

func main() {
	loggingSettings("shares.log")

	log.Println("Webサーバーを開始します...")
	r := newRoom()
	http.Handle("/room", r)
	go r.run()
	server := http.Server{
		Addr: ":9000",
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/home", templateHandler)
	http.HandleFunc("/stickies", getStickiesInfo)
	http.HandleFunc("/load-sticky-id", loadStickyId)
	http.HandleFunc("/create-sticky", createSticky)
	http.HandleFunc(("/update-sticky"), updateSticky)
	server.ListenAndServe()
}
