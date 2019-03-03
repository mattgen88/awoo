package handler

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/Sigafoos/awoo/game"
	"github.com/Sigafoos/awoo/player"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Disable CORS for testing
	},
}

type Handler struct {
	game     *game.Game
	joinChan chan *player.Player
	ips      map[string]*player.Player
}

func New() *Handler {
	joinChan := make(chan *player.Player)
	return &Handler{
		game:     game.New(joinChan),
		joinChan: joinChan,
		ips:      make(map[string]*player.Player),
	}
}

func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("BAD upgrade:", err)
		return
	}

	// handle reconnects
	ip := h.ipFromRequest(r)
	if p, exists := h.ips[ip]; exists {
		p.Reconnect(c)
	} else {
		p := player.New(c, h.joinChan)
		h.ips[ip] = p
		go p.Play()
	}
}

// Awoo's only function is to amuse its authors.
func (h *Handler) Awoo(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v awoos\n", r.RemoteAddr)
	fmt.Fprintln(w, "awoooooooooooo")
}

func (h *Handler) ipFromRequest(r *http.Request) string {
	// This bit hasn't been tested. It also doesn't handle case differences.
	ips := r.Header.Get("X-Forwarded-For")
	if ips != "" {
		list := strings.Split(ips, ", ")
		return list[0]
	}

	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return ""
}
