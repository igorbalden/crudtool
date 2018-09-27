package main

/**
 * Crudtool
 * License: MIT
 **/
import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/igorbalden/crudtool/config"
	"github.com/igorbalden/crudtool/middles"
	"github.com/igorbalden/crudtool/views"
)

func main() {
	values, _ := config.ReadConfig("./config/config.json")
	var port *string

	// port number from cli -port flag
	port = flag.String("port", "", "IP address")
	flag.Parse()

	if *port != "" {
		//We expect :8080 as input, else append ':'
		if !strings.HasPrefix(*port, ":") {
			*port = ":" + *port
		}
		values.ServerPort = *port
	}

	if values.ServerPort == "" {
		log.Print(`A server port number is needed in config.json 
			or as a command line flag "-port" `)
		os.Exit(1)
	}

	//Load all Template files
	views.PopulateTemplates()

	// Routes
	rt := mux.NewRouter()
	rt = rt.StrictSlash(true)
	rt.HandleFunc("/login", views.LoginFunc)
	rt.HandleFunc("/logout", auth.RequiresLogin(views.LogoutFunc))
	rt.HandleFunc("/dbtable/{dbtable}/dbname/{dbname}/p/{p}", auth.RequiresLogin(views.TblContent))
	rt.HandleFunc("/dbtable/{dbtable}/dbname/{dbname}", auth.RequiresLogin(views.TblContent))
	rt.HandleFunc("/dbname/{dbname}/p/{p}", auth.RequiresLogin(views.ListTables))
	rt.HandleFunc("/dbname/{dbname}", auth.RequiresLogin(views.ListTables))
	rt.HandleFunc("/sqlstmt", auth.RequiresLogin(views.SQLStmt))
	rt.HandleFunc("/dbselect", auth.RequiresLogin(views.DbSelect))
	rt.HandleFunc("/dbselect/p/{p}", auth.RequiresLogin(views.DbSelect))
	rt.HandleFunc("/", auth.RequiresLogin(views.DbSelect))

	rt.PathPrefix("/static/").Handler(http.FileServer(http.Dir("public")))
	rt.Handle("/favicon.ico", http.FileServer(http.Dir("public")))

	// Serve
	wait, _ := time.ParseDuration(values.Wait)
	wTime, _ := time.ParseDuration(values.WriteTimeout)
	rTime, _ := time.ParseDuration(values.ReadTimeout)
	iTime, _ := time.ParseDuration(values.IdleTimeout)
	srv := &http.Server{
		Addr: values.ServerPort,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: wTime,
		ReadTimeout:  rTime,
		IdleTimeout:  iTime,
		Handler:      rt, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Print("Running on port ", values.ServerPort)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
