package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/travissimon/house/colour"
	"github.com/travissimon/huego"
	"log"
	"math/rand"
	"net/http"
	"strconv"
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
				s.clients[c.id] = c
				// delete a client
			case c := <-s.deleteChannel:
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
	case "setGenerator":
		s.HandleSetGenerator(msg)
	case "saveScene":
		s.HandleSaveScene(msg)
	case "setScene":
		s.HandleSetScene(msg)
	case "deleteScene":
		s.HandleDeleteScene(msg)
	case "setScheme":
		s.HandleSetScheme(msg)
	case "saveScheme":
		s.HandleSaveScheme(msg)
	case "deleteScheme":
		s.HandleDeleteScheme(msg)
	case "allOn":
		s.HandleAllOn(msg)
	case "allOff":
		s.HandleAllOff(msg)
	case "setPower":
		s.HandleSetPower(msg)
	default:
		log.Printf("Unknown action: '%v'\n", msg.Request.Action)
	}
}

func stopSceneActivity() {
	wemoUtil.Stop()
	mainScene.Stop()
}

func startSceneActivity(scene *Scene) {
	log.Printf("Starting scene activity: %s", scene.Name)
	mainScene = scene
	mainScene.Start()
	wemoUtil.Start(notifySceneOfActivity)
}

func (s *SocketServer) HandleSetLightColour(msg SocketMessage) {
	stopSceneActivity()
	args := msg.GetSetLightArguments()
	l := getLightById(args.Id)
	l.SetColourFromHex(args.Hex)
	l.SetState()
}

func (s *SocketServer) HandleSaveScene(msg SocketMessage) {
	args := msg.GetSaveSceneArguments()
	scene := NewScene()
	scene.Id = args.Id
	scene.Name = args.Name
	scene.ActiveTransition = args.ActiveTransition
	scene.ActiveHold = args.ActiveHold
	scene.InactiveTransition = args.InactiveTransition
	scene.InactiveHold = args.InactiveHold
	scene.ActiveScheme = args.ActiveScheme
	scene.InactiveSchemes = args.InactiveSchemes
	scene.Persist()
}

func (s *SocketServer) HandleSetScene(msg SocketMessage) {
	log.Printf("Set scene called. Stopping any current activity")
	stopSceneActivity()
	args := msg.GetSetSceneArguments()
	log.Printf("Loading scene: %v\n", args.Id)
	sc, _ := LoadSceneById(args.Id)
	startSceneActivity(sc)
}

func (s *SocketServer) HandleDeleteScene(msg SocketMessage) {
	args := msg.GetDeleteSceneArguments()
	DeleteSceneById(args.Id)
}

func (s *SocketServer) HandleSetGenerator(msg SocketMessage) {
	stopSceneActivity()
	args := msg.GetSetGeneratorArguments()
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
		time.Sleep(30 * time.Millisecond)
	}

	proxies := s.getCurrentLightProxies()
	msg.Response = proxies
	s.sendAll(msg)
}

func (s *SocketServer) HandleSetScheme(msg SocketMessage) {
	stopSceneActivity()
	args := msg.GetSetSchemeArguments()
	scheme, _ := LoadSchemeById(args.Id)
	for _, light := range scheme.Lights {
		l := getLightById(light.Id)
		l.SetColourFromHex(light.Hex)
		l.SetStateWithTransition(10)
		time.Sleep(30 * time.Millisecond)
	}
}

func (s *SocketServer) HandleSaveScheme(msg SocketMessage) {
	args := msg.GetSaveSchemeArguments()
	scheme := NewScheme()
	if args.Id != "" {
		scheme.Id, _ = strconv.Atoi(args.Id)
	}
	scheme.Name = args.Name
	scheme.Lights = args.Lights
	scheme.Persist()
}

func (s *SocketServer) HandleDeleteScheme(msg SocketMessage) {
	args := msg.GetDeleteSchemeArguments()
	DeleteSchemeById(args.Id)
}

func (s *SocketServer) HandleSetPower(msg SocketMessage) {
	stopSceneActivity()
	args := msg.GetSetPowerArguments()
	l := getLightById(args.Id)
	if l.State.On != args.TurnOn {
		l.State.On = args.TurnOn
		l.SetState()
	}
}

func (s *SocketServer) HandleAllOn(msg SocketMessage) {
	stopSceneActivity()
	powerAll(true)
}

func (s *SocketServer) HandleAllOff(msg SocketMessage) {
	stopSceneActivity()
	powerAll(false)
}

func powerAll(turnOn bool) {
	stopSceneActivity()
	for _, light := range lights {
		if light.State.On != turnOn {
			light.State.On = turnOn
			if turnOn {
				light.SetStateWithTransition(10)
			} else {
				light.SetState()
			}
		}
		time.Sleep(30 * time.Millisecond)
	}
}

func (s *SocketServer) getCurrentLightProxies() []*LightProxy {
	proxies := make([]*LightProxy, 0, len(lights))
	for _, light := range lights {
		proxy := NewLightProxyFromLight(light)
		proxies = append(proxies, proxy)
	}
	return proxies
}
