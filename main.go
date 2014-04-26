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

var templates = template.Must(template.ParseFiles("views/home.html", "views/lights.html", "views/generator.html", "views/scheme_index.html"))

func homeHandler(w http.ResponseWriter, r *http.Request) {
	proxies := getLightProxies()
	templates.ExecuteTemplate(w, "home.html", proxies)
}

func lightsHandler(w http.ResponseWriter, r *http.Request) {
	reloadLightState()
	proxies := getLightProxies()
	templates.ExecuteTemplate(w, "lights.html", proxies)
}

func randomSchemeHandler(w http.ResponseWriter, r *http.Request) {
	reloadLightState()
	proxies := getLightProxies()
	templates.ExecuteTemplate(w, "generator.html", proxies)
}

func schemeHandler(w http.ResponseWriter, r *http.Request) {
	schemes, _ := LoadSchemes()
	templates.ExecuteTemplate(w, "scheme_index.html", schemes)
}

func getLightProxies() []*LightProxy {
	proxies := make([]*LightProxy, 0, len(lights))
	for _, light := range lights {
		proxies = append(proxies, NewLightProxyFromLight(light))
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

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/lights", lightsHandler)
	http.HandleFunc("/generator", randomSchemeHandler)
	http.HandleFunc("/schemes", schemeHandler)

	log.Printf("Listening on port %s\n", port)
	http.ListenAndServe(port, nil)
}

func reloadLightState() []*huego.Light {
	lights, _ = base.GetLights()
	log.Printf("Reloaded %v lights\n", len(lights))
	return lights
}

func initHue() {
	bases, _ := huego.DiscoverBases()
	base = &bases[0]
	lights, _ = base.GetLights()
}
