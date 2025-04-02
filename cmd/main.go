package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/thepralad/socialnetwork/internal/handlers"
	"github.com/thepralad/socialnetwork/internal/models"
	"github.com/thepralad/socialnetwork/pkg/render"
)

func main() {
	db, err := sql.Open("mysql", "avnadmin:AVNS_0gwggs0ttt3MkMqPLEJ@tcp(mysql-152cca74-snet.i.aivencloud.com:21979)/snet")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize templates
	if err := render.Init(); err != nil {
		log.Fatal("Failed to initialize templates:", err)
	}

	store := models.MySQLUserStore{DB: db}

	authHandler := handlers.NewAuthHandler(&store)
	postHandler := handlers.NewPostHandler(&store)
	userHandler := handlers.NewUserHandler(&store)

	mux := http.NewServeMux()

	//Serve static files
	fs := http.FileServer(http.Dir("internal/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", authHandler.HomeHandler)
	mux.HandleFunc("/register", authHandler.RegisterHandler)
	mux.HandleFunc("/login", authHandler.LoginHandler)
	mux.HandleFunc("/verify", authHandler.VerifyHandler)
	mux.HandleFunc("/logout", authHandler.LogoutHandler)

	mux.HandleFunc("/home", postHandler.HomePostHandler)
	mux.HandleFunc("/home/", postHandler.UserProfileHandler)
	mux.HandleFunc("/editprofile", postHandler.EditProfileHandler)
	mux.HandleFunc("/getposts", postHandler.GetPostHandler)
	mux.HandleFunc("/post", postHandler.PostHandler)
	mux.HandleFunc("/updatemetric", postHandler.UpdateMetricHandler)

	mux.HandleFunc("/poke", userHandler.PokeHandler)
	mux.HandleFunc("/pokes", userHandler.PokesHandler)

	log.Print("starting server at :8080")
	http.ListenAndServe(":8080", mux)
}
