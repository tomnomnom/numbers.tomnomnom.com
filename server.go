package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

var templates = template.Must(template.ParseFiles("main.html"))
var state = struct {
	sync.Mutex
	mean int
}{mean: 500000000}

type Client struct {
	mean chan string
}

var newClients = make(chan Client)
var deadClients = make(chan Client)

func main() {
	log.Println("Starting...")

	rand.Seed(time.Now().UnixNano())

	// Load the page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check for a clientId cookie
		clientId := fmt.Sprintf("%x%x", rand.Int(), rand.Int())
		c, err := r.Cookie("clientId")
		if err == nil {
			clientId = c.Value
		}

		state.Lock()
		mean := state.mean
		state.Unlock()

		http.SetCookie(w, &http.Cookie{Name: "clientId", Value: clientId})
		err = templates.ExecuteTemplate(w, "main.html", struct {
			ClientId string
			Mean     int
		}{clientId, mean})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Handle new guesses
	newGuesses := make(chan string)
	http.Handle("/answers", websocket.Handler(func(ws *websocket.Conn) {
		c := Client{mean: make(chan string)}
		newClients <- c

		// Send new means when we get them
		go func() {
			for {
				mean := <-c.mean
				ws.Write([]byte(mean))
			}
		}()

		// Wait for guesses
		s := bufio.NewScanner(ws)
		for s.Scan() {
			msg := s.Text()
			log.Printf("New guess: %s\n", msg)
			newGuesses <- msg

			if err := s.Err(); err != nil {
				log.Println("Client error: ", err)
				break
			}
		}

		deadClients <- c
	}))

	http.Handle("/static/", http.FileServer(http.Dir(".")))
	go http.ListenAndServe(":8080", nil)

	// Meat
	clients := make(map[Client]interface{})
	guesses := make(map[string]int)
	for {
		select {

		case client := <-newClients:
			clients[client] = struct{}{}

		case client := <-deadClients:
			delete(clients, client)

		case guess := <-newGuesses:
			// Calculate the mean and send out
			parts := strings.SplitN(guess, ":", 2)
			id := parts[0]
			val, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}
			guesses[id] = val

			t := 0
			for _, g := range guesses {
				t += g
			}
			mean := t / len(guesses)
			state.Lock()
			state.mean = mean
			state.Unlock()

			for client, _ := range clients {
				client.mean <- strconv.Itoa(mean)
			}

		}
	}

}
