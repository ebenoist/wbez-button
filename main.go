package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"sync"
	"syscall"
)

type Player struct {
	cmd *exec.Cmd
	sync.Mutex
}

func (p *Player) isPlaying() bool {
	return p.cmd != nil
}

func (p *Player) Play(url string) error {
	p.Lock()
	defer p.Unlock()

	if !p.isPlaying() {
		p.cmd = exec.Command("mplayer", url)
		return p.cmd.Start()
	}

	return nil
}

func (p *Player) Stop() error {
	p.Lock()
	defer p.Unlock()

	if p.isPlaying() {
		p.cmd.Process.Signal(syscall.SIGTERM)
		p.cmd.Wait()
		p.cmd = nil
		return nil
	}

	return nil
}

func main() {
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)
	player := &Player{}

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		if player.isPlaying() {
			w.Write([]byte("streaming"))
		} else {
			w.Write([]byte("stopped"))
		}
	})

	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		err := player.Play("http://stream.wbez.org/wbez128.mp3")
		if err != nil {
			http.Error(w, fmt.Sprintf("could not start: %s", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("streaming"))
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		err := player.Stop()
		if err != nil {
			http.Error(w, fmt.Sprintf("not playing: %s", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("stopped"))
	})

	http.ListenAndServe(":8080", nil)
}
