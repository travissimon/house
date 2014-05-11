package main

import (
	"github.com/travissimon/house/colour"
	"github.com/travissimon/huego"
	"html/template"
	"log"
	"net/http"
)

var base *huego.Base
var lights []*huego.Light
var generator *colour.ColourSchemeGenerator = colour.NewColourSchemeGenerator()
var mainScene *Scene
var wemoUtil = NewWemoUtil()

func notifySceneOfActivity(sensorName string) {
	if mainScene != nil {
		mainScene.NotifyOfMovement(sensorName)
	}
}

func modIsZero(val int, modAmount int, expValue int) bool {
	retVal := val%modAmount == expValue
	return retVal
}

func isSelectedOption(selectedOptions []int, currentOption int) bool {
	for _, id := range selectedOptions {
		if id == currentOption {
			return true
		}
	}
	return false
}

var funcMap = template.FuncMap{
	"modIsZero":        modIsZero,
	"isSelectedOption": isSelectedOption,
}

var templates = template.Must(template.New("templates").Funcs(funcMap).ParseGlob("views/*.html"))

func homeHandler(w http.ResponseWriter, r *http.Request) {
	reloadLightState()
	proxies := getLightProxies()
	templates.ExecuteTemplate(w, "home", proxies)
}

func editSchemeHandler(w http.ResponseWriter, r *http.Request) {
	reloadLightState()
	id := r.URL.Path[len("/schemes/edit/"):]
	scheme, _ := LoadSchemeById(id)
	templates.ExecuteTemplate(w, "lights", scheme)
}

func deleteSchemeHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/schemes/delete/"):]
	scheme, _ := LoadSchemeById(id)
	templates.ExecuteTemplate(w, "deleteScheme", scheme)
}

func lightsHandler(w http.ResponseWriter, r *http.Request) {
	reloadLightState()
	scheme := NewScheme()
	scheme.Lights = getLightArgs()
	templates.ExecuteTemplate(w, "lights", scheme)
}

func randomSchemeHandler(w http.ResponseWriter, r *http.Request) {
	reloadLightState()
	proxies := getLightProxies()
	templates.ExecuteTemplate(w, "generator", proxies)
}

type ScenePage struct {
	Scene     *Scene
	Schemes   []*Scheme
	AllScenes []*Scene
}

func schemeHandler(w http.ResponseWriter, r *http.Request) {
	schemes, _ := LoadSchemes()
	templates.ExecuteTemplate(w, "schemeIndex", schemes)
}

func sceneHandler(w http.ResponseWriter, r *http.Request) {
	scenes, _ := LoadScenes()
	pg := &ScenePage{nil, nil, scenes}
	templates.ExecuteTemplate(w, "sceneIndex", pg)
}

func editSceneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/scenes/edit/"):]
	schemes, _ := LoadSchemes()
	scenes, _ := LoadScenes()
	scene, _ := LoadSceneById(id)
	if scene == nil {
		scene = NewScene()
	}
	pg := &ScenePage{scene, schemes, scenes}
	templates.ExecuteTemplate(w, "sceneEdit", pg)
}

func getLightProxies() []*LightProxy {
	proxies := make([]*LightProxy, 0, len(lights))
	for _, light := range lights {
		proxies = append(proxies, NewLightProxyFromLight(light))
	}
	return proxies
}

func getLightArgs() []*SetLightArguments {
	proxies := make([]*SetLightArguments, 0, len(lights))
	for _, light := range lights {
		proxies = append(proxies, &SetLightArguments{light.Id, light.Name, light.ToHex()})
	}
	return proxies
}

func main() {
	initHue()
	port := ":8080"

	// websocket server
	socketServer := NewSocketServer()
	socketServer.Listen()

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("static/images"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("static/fonts"))))

	http.HandleFunc("/lights", lightsHandler)
	http.HandleFunc("/generator", randomSchemeHandler)
	http.HandleFunc("/schemes/edit/", editSchemeHandler)
	http.HandleFunc("/schemes/delete/", deleteSchemeHandler)
	http.HandleFunc("/schemes", schemeHandler)
	http.HandleFunc("/scenes/edit/", editSceneHandler)
	http.HandleFunc("/scenes", sceneHandler)
	http.HandleFunc("/", homeHandler)

	log.Printf("Listening on port %s\n", port)
	http.ListenAndServe(port, nil)
}

func reloadLightState() []*huego.Light {
	lights, _ = base.GetLights()
	return lights
}

func initHue() {
	var err error
	base, err = GetHueBase()
	if err != nil {
		log.Printf("Error: %v", err)
	}
	log.Printf("init base: %v\n", base)
	//bases, _ := huego.DiscoverBases()
	//base = &bases[0]
	log.Printf("Base ip: %v, developer: %v", base.InternalIp, base.Username)
	lights, _ = base.GetLights()
}
