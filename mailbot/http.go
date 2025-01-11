package mailbot

import (
	"context"
	"io"
	"log"
	"mail2telegram/env"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

func (bot *Bot) StartHttpServer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	mux := http.NewServeMux()
	mux.HandleFunc("/"+env.Env.StoragePrefix+"/{id}", getMail)

	server := &http.Server{Addr: env.Env.HTTPAddr, Handler: mux}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Listen&Serve result: %e\n", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(shutdownCtx)
	if err != nil {
		log.Printf("Shutdown error: %e\n", err)
	}
	log.Println("Http server shutted down")
}

func getMail(writer http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Invalid UUID: %e", err)
		io.WriteString(writer, "Invalid UUID")
		return
	}

	file, err := OpenFile(uid)
	if err != nil {
		log.Printf("Error while reading file: %e", err)
		io.WriteString(writer, "File not found")
		return
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		log.Printf("Error while reading file: %e", err)
		io.WriteString(writer, "File read error")
		return
	}
}
