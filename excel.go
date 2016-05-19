package main

import (
	"fmt"
	"net/http"
	"database/sql"
	"github.com/tealeg/xlsx"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func main(){
	http.HandleFunc("/", handle)
	http.ListenAndServe(":8080", nil)
}

func handle(w http.ResponseWriter, r *http.Request){
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error
	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("mysql", "root:Tutululu@tcp(192.168.1.132:3306)/wealthbetter_com")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM `wb_sina_log`")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err)
		}
		row = sheet.AddRow()
		for _, cellValue := range values {
			cell = row.AddCell()
			cell.Value = string(cellValue)
		}
	}
	w.Header().Set("Content-Type", "application/vnd.ms-excel;charset=UTF-8")
	w.Header().Set("Pragma", "public")
	w.Header().Set("Expires", "0")
	w.Header().Set("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Content-Type", "application/force-download")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Type", "application/download")
	w.Header().Set("Content-Disposition", "attachment;filename="+time.Now().Format("2006-01-02 15:04:05")+".xlsx")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	err = file.Write(w)
	if err != nil {
		fmt.Printf(err.Error())
	}
}
