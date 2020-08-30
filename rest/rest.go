package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gorilla/mux"
)

var excelFile string
var port int

func StartServer(f string, p int) {
	excelFile = f
	port = p
	router := mux.NewRouter()
	router.HandleFunc("/", getSheetNames)
	router.HandleFunc("/{sheet}", getSheetData)

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("localhost:%d", port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Printf("http://localhost:%d\n", port)
	log.Fatal(srv.ListenAndServe())
}

func getSheetNames(w http.ResponseWriter, r *http.Request) {
	f, err := excelize.OpenFile(excelFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	sheetMap := f.GetSheetMap()
	sheetNames := make([]string, 0, len(sheetMap))
	for _, v := range sheetMap {
		v = strings.ReplaceAll(v, " ", "%20")
		sheetNames = append(sheetNames, v)
	}
	json.NewEncoder(w).Encode(sheetNames)
}

func getSheetData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sheet := vars["sheet"]
	file, err := excelize.OpenFile(excelFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Prepare data
	data := file.GetRows(sheet)

	json.NewEncoder(w).Encode(data)
}
