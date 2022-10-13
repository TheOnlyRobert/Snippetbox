package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"rfranks/snippetbox/pkg/models/mysql"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

// The application struct holds the application-wide dependencies for the web application.
type application struct {
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
	errorLog      *log.Logger
	infoLog       *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySql data source")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		errorLog:      errorLog,
		infoLog:       infoLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
