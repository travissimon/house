package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/travissimon/house/colour"
	"github.com/travissimon/huego"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Socket Server
type SocketServer struct {
	clients        map[int]*SocketClient
	addChannel     chan *SocketClient
	deleteChannel  chan *SocketClient
	sendAllChannel chan SocketMessage
	doneChannel    chan bool
	errorChannel   chan error
}

func NewSocketServer() *SocketServer {
	rand.Seed(time.Now().UnixNano())

	clients := make(map[int]*SocketClient)
	addChannel := make(chan *SocketClient)
	deleteChannel := make(chan *SocketClient)
	sendAllChannel := make(chan SocketMessage)
	doneChannel := make(chan bool)
	errorChannel := make(chan error)

	return &SocketServer{
		clients,
		addChannel,
		deleteChannel,
		sendAllChannel,
		doneChannel,
		errorChannel,
	}
}

func (s *SocketServer) Add(c *SocketClient) {
	s.addChannel <- c
}

func (s *SocketServer) Delete(c *SocketClient) {
	s.deleteChannel <- c
}

func (s *SocketServer) SendAll(msg SocketMessage) {
	s.sendAllChannel <- msg
}

func (s *SocketServer) Err(err error) {
	s.errorChannel <- err
}

func (s *SocketServer) sendAll(msg SocketMessage) {
	for _, c := range s.clients {
		c.Write(msg)
	}
}

// Listens and broadcasts client events
func (s *SocketServer) Listen() {

	// websocket handler func
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errorChannel <- err
			}
		}()

		client := NewClient(ws, s)
		s.Add(client)
		client.Listen()
	}

	http.Handle("/socket", websocket.Handler(onConnected))
	log.Println("Created Socket Handler on: /socket")

	go func() {
		for {
			select {

			// add new client
			case c := <-s.addChannel:
				log.Println("Adding new client: ", c.id)
				s.clients[c.id] = c

				// delete a client
			case c := <-s.deleteChannel:
				log.Println("Deleting client: ", c.id)
				delete(s.clients, c.id)

				// broadcast a message
			case msg := <-s.sendAllChannel:
				s.sendAll(msg)

			case <-s.doneChannel:
				return
			}
		}
	}()
}

func getLightById(id string) *huego.Light {
	for _, l := range lights {
		if l.Id == id {
			return l
		}
	}
	return nil
}

func (s *SocketServer) HandleIncomingRequest(msg SocketMessage) {
	switch msg.Request.Action {
	case "setLightColour":
		s.HandleSetLightColour(msg)
	case "setScheme":
		s.HandleSetScheme(msg)
	default:
		log.Printf("Unknown action: '%v'\n", msg.Request.Action)
	}
}

func (s *SocketServer) HandleSetLightColour(msg SocketMessage) {
	args := msg.GetSetLightArguments()
	l := getLightById(args.Id)
	l.SetColourFromHex(args.Value)
	l.SetState()
}

func (s *SocketServer) HandleSetScheme(msg SocketMessage) {
	args := msg.GetSetSchemeArguments()
	strategy := colour.GetHarmonyStrategy(args.Strategy)
	generator.SetStrategy(strategy)
	generator.SetAngle(args.Angle)
	generator.SetTint(args.Tint)
	generator.SetShade(args.Shade)

	primaryColour, _ := colour.FromHex(args.PrimaryColour)
	schemeColours := generator.GetScheme(primaryColour)

	if len(schemeColours) == 0 {
		log.Printf("Error generating scheme: %v, %v, %v, %v", strategy, args.Angle, args.Tint, args.Shade)
		return
	}

	colours := make([]*colour.Colour, 0, len(schemeColours)*5)
	for _, sc := range schemeColours {
		colours = append(colours, &sc.Colour)
		colours = append(colours, &sc.Tints[0])
		colours = append(colours, &sc.Tints[1])
		colours = append(colours, &sc.Shades[0])
		colours = append(colours, &sc.Shades[1])
	}

	for _, light := range lights {
		index := rand.Int31n(int32(len(colours)))
		c := colours[index]
		h, s, v := c.ToHue()
		light.State.Hue = h
		light.State.Sat = s
		light.State.Bri = v
		light.SetState()
	}

	proxies := s.getCurrentLightProxies()
	msg.Response = proxies
	s.sendAll(msg)
}

func (s *SocketServer) getCurrentLightProxies() []*LightProxy {
	proxies := make([]*LightProxy, 0, len(lights))
	for _, light := range lights {
		proxy := NewLightProxyFromLight(light)
		proxies = append(proxies, proxy)
	}
	return proxies
}
