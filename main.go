package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"syscall"
)

var cmd = exec.Command("mplayer", "http://stream.wbez.org/wbez128.mp3")

func main() {
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		err := cmd.Start()
		if err != nil {
			http.Error(w, fmt.Sprintf("could not start: %s", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("streaming"))
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		if cmd.Process != nil {
			cmd.Process.Signal(syscall.SIGTERM)
			cmd.Wait()
			cmd = exec.Command("mplayer", "http://stream.wbez.org/wbez128.mp3")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("stopped"))
			return
		}

		http.Error(w, fmt.Sprintf("not playing"), http.StatusBadRequest)
	})

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}
