package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/nopdan/rose/pkg/core"
)

type Config struct {
	Name    string `json:"name"`
	IFormat string `json:"inputFormat"`
	Kind    string `json:"kind"`
	Schema  string `json:"schema"`
	Rule    string `json:"rule"`
	OFormat string `json:"outputFormat"`
}

//go:embed dist
var dist embed.FS

func Serve(port int) {
	dist, _ := fs.Sub(dist, "dist")
	http.Handle("/", http.FileServer(http.FS(dist)))

	var iData []byte
	var mbData []byte
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w)
		if r.Method != "POST" {
			return
		}
		if r.Body == nil {
			log.Printf("POST:/api no request body\n")
			fmt.Fprint(w, "error")
			return
		}

		decoder := json.NewDecoder(r.Body)
		c := Config{}
		decoder.Decode(&c)
		log.Printf("POST:/api request Body: %+v\n", c)

		var conf = &core.Config{}
		conf.IName = c.Name
		conf.IFormat = c.IFormat
		if iData == nil {
			log.Printf("POST:/api no input data\n")
			fmt.Fprint(w, "error")
			return
		}
		conf.IData = iData
		conf.Schema = c.Schema
		conf.MbData = mbData
		if c.Rule == "AABC" {
			conf.AABC = true
		}
		conf.OFormat = c.OFormat
		data := conf.Marshal()
		conf.Save(data)
		log.Printf("POST:/api 保存到: %v\n\n\n", conf.OName)

		fmt.Fprint(w, conf.OName)
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w)
		log.Println("GET:/list")
		b, _ := json.Marshal(core.FormatList)
		w.Write(b)
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		iData = uploadHandler(&w, r, "/upload")
	})

	http.HandleFunc("/upload/mb", func(w http.ResponseWriter, r *http.Request) {
		mbData = uploadHandler(&w, r, "/upload/mb")
	})

	sport := fmt.Sprint(port)
	log.Println("Listening on http://localhost:" + sport)
	openBrowser("http://localhost:" + sport)
	http.ListenAndServe(":"+sport, nil)
}

func uploadHandler(w *http.ResponseWriter, r *http.Request, path string) []byte {
	setupCORS(w)
	if r.Method != "POST" {
		return nil
	}
	// 最大 1 GB
	r.ParseMultipartForm(1 << 32)
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("POST:%s err: %v\n", path, err)
		return nil
	}
	defer file.Close()
	log.Printf("POST:%s %v", path, handler.Header)

	mbData, err := io.ReadAll(file)
	if err != nil {
		log.Printf("POST:%s err: %v\n", path, err)
		return nil
	}
	fmt.Fprintf((*w), "%v", handler.Header)
	return mbData
}

func setupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func openBrowser(url string) {
	var name string
	switch runtime.GOOS {
	case "windows":
		name = "explorer"
	case "linux":
		name = "xdg-open"
	default:
		name = "open"
	}
	cmd := exec.Command(name, url)
	cmd.Start()
}
