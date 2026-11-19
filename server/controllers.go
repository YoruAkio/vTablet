package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Teages/vTablet/internal/pointer"
	"golang.org/x/net/websocket"
)

var (
	httpPort = flag.Int("port", 23052, "HTTP & WebSocket service port")
)

func initWebServices(port int) {
	if *httpPort != 23052 {
		port = *httpPort
	}

	http.HandleFunc("/", clientPage)
	http.HandleFunc("/connect", connecter)

	for _, screen := range pointer.GetScreens() {
		http.Handle(fmt.Sprintf("/digi/%s", screen.Uid), connectFactory(screen.Uid))
	}
	http.HandleFunc("/digi/", func(w http.ResponseWriter, r *http.Request) {
		uid := strings.TrimPrefix(r.URL.Path, "/digi/")
		connectFactory(uid).ServeHTTP(w, r)
	})

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("Starting vTablet web service on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func connecter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	screens := pointer.GetScreens()
	jsonStr, _ := json.Marshal(screens)
	io.WriteString(w, string(jsonStr))
}

func connectFactory(screeUid string) websocket.Handler {
	return websocket.Handler(func(c *websocket.Conn) {
		if c.Config().Origin != nil {
			origin := c.Config().Origin.Host
			allows := "http://localhost:23052\n" + "https://vtablet.yoruakio.xyz\n" + "https://vtablet.teages.xyz\n"
			if origin != "" && !strings.Contains(allows, origin) &&
				!strings.HasPrefix(origin, "192.168.") &&
				!strings.HasPrefix(origin, "10.") &&
				!strings.HasPrefix(origin, "172.") &&
				!strings.HasPrefix(origin, "localhost") &&
				!strings.HasPrefix(origin, "127.0.0.1") {
				c.Close()
				return
			}
		}
		pointerServices(screeUid, c)
	})
}

func clientPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://vtablet.yoruakio.xyz", http.StatusFound)
}
