package formosa

import (
	_ "crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func WebServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/status", StatusHandler)
	http.Handle("/", r)
	WsServer()
	log.Println("Web Service started.")
	http.ListenAndServe(fmt.Sprintf("%s:%d", CONFIGS.Web.Host, CONFIGS.Web.Port), nil)
}

func jsonParser(data interface{}, w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if data != nil {
		json, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(500)
			log.Println("Error generating json", err)
			fmt.Fprintln(w, "Could not generate JSON")
			return
		}
		fmt.Fprint(w, string(json))
	} else {
		w.WriteHeader(404)
		fmt.Fprint(w, "404 no data can be find.")
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(403)
	fmt.Fprint(w, "403 Forbidden")
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	data["connection"] = fmt.Sprintf("%d", ProxyConn)
	jsonParser(data, w)
}
