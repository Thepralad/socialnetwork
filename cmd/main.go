package main

import (
	"log"
	"net/http"
	"github.com/thepralad/socialnetwork/internal/handlers"
	"github.com/thepralad/socialnetwork/internal/models"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main(){
	db, err := sql.Open("mysql", "avnadmin:AVNS_0gwggs0ttt3MkMqPLEJ@tcp(mysql-152cca74-snet.i.aivencloud.com:21979)/snet")	
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()
	
	store := models.MySQLUserStore{DB: db}
	handler := handlers.NewUserHandler(&store)

	mux := http.NewServeMux()
	
	mux.HandleFunc("/", handler.HomeHandler)
	mux.HandleFunc("/register", handler.RegisterUser)
	
	log.Print("starting server at :8080")
	http.ListenAndServe(":8080", mux)
}
