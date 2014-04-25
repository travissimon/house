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

var templates = template.Must(template.ParseFiles("views/home.html", "views/scheme.html"))

func homeHandler(w http.ResponseWriter, r *http.Request) {
	proxies := make([]*LightProxy, 0, len(lights))
	for _, light := range lights {
		proxy := NewLightProxyFromLight(light)
		proxies = append(proxies, proxy)
	}
	templates.ExecuteTemplate(w, "home.html", proxies)
}

func schemeHandler(w http.ResponseWriter, r *http.Request) {
	proxies := make([]*LightProxy, 0, len(lights))
	for _, light := range lights {
		proxy := NewLightProxyFromLight(light)
		proxies = append(proxies, proxy)
	}
	templates.ExecuteTemplate(w, "scheme.html", proxies)
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
	http.HandleFunc("/scheme", schemeHandler)

	log.Printf("Listening on port %s\n", port)
	http.ListenAndServe(port, nil)

}

func initHue() {
	bases, _ := huego.DiscoverBases()
	base = &bases[0]
	lights, _ = base.GetLights()
}
