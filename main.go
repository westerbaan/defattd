package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/timshannon/bolthold"
)

type EMailRecord struct {
	EMail string
}

var (
	db *bolthold.Store
)

func handleExport(w http.ResponseWriter, r *http.Request) {
	var rrs []EMailRecord
	db.Find(&rrs, nil)
	f, err := os.Create("defattd-email.txt")
	if err != nil {
		log.Printf("export os.Create(): %v", err)
		return
	}
	defer f.Close()
	for _, rr := range rrs {
		f.WriteString(fmt.Sprintf("%s\n", rr.EMail))
	}
}

func handleEMail(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Expected POST", 500)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("ParseForm() err: %v", err), 500)
		return
	}
	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Missing email field", 500)
		return
	}
	db.Insert(bolthold.NextSequence(), EMailRecord{
		EMail: r.FormValue("email"),
	})
}

func main() {
	var bindAddr string

	flag.StringVar(&bindAddr, "bind", "localhost:1405", "address to bind to")
	flag.Parse()

	http.HandleFunc("/email", handleEMail)
	http.HandleFunc("/export", handleExport)

	var err error
	db, err = bolthold.Open("defattd.db", 0600, nil)
	if err != nil {
		log.Fatalf("bolthold.Open(): %v", err)
	}

	log.Fatal(http.ListenAndServe(bindAddr, nil))
}
