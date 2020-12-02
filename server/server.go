package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"lrucached/cache"
	"net/http"
	"time"
)

const (
	mainURL = "/api"
	version = "v1"
)

// CacheServer implementation cache server on http protocols
type CacheServer struct {
	mux       *http.ServeMux
	c         cache.Cache
	indexPath string
	statPath  string
}

// NewCacheServer -creates new instance of cache server
func NewCacheServer(capacity uint64, ttl, ttc time.Duration) *CacheServer {
	return &CacheServer{
		mux:       http.NewServeMux(),
		c:         cache.NewCache(capacity, ttl, ttc),
		indexPath: fmt.Sprintf("%s/%s/%s", mainURL, version, "cached"),
		statPath:  fmt.Sprintf("%s/%s/%s", mainURL, version, "stat"),
	}
}

// Init -creates lrucached instance and initializes routes and configs
func (cs *CacheServer) Init() {

	// Read configs from os.Env

	// Adding handlers
	println("above line of handler")
	cs.mux.Handle(cs.indexPath, cs.indexHandler())
	// cs.mux.Handle(cs.statPath, cs.statHandler())
}

// Run launches the server on specified port
func (cs *CacheServer) Run(addr string) {
	fmt.Printf("Lrucached is now running on %s address\n", addr)
	log.Fatal(http.ListenAndServe(addr, cs.mux))
}

func (cs *CacheServer) indexHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		switch r.Method {
		case http.MethodPut:
			cs.setHandler(w, r)
		case http.MethodGet:
			cs.getHandler(w, r)
		case http.MethodDelete:
			cs.removeHandler(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
	return http.HandlerFunc(fn)
}

func (cs *CacheServer) setHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len(cs.indexPath):]
	if key == "" {
		responseError(w, http.StatusBadRequest, "empty key passed")
		return
	}
	rawData, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		// add response
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	cs.c.Set(key, rawData)
	w.WriteHeader(http.StatusCreated)
}

func (cs *CacheServer) getHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len(cs.indexPath):]
	log.Println(r.URL.Path)
	if key == "" {
		responseError(w, http.StatusBadRequest, "empty key passed")
		return
	}
	rawData, found := cs.c.Get(key)
	if !found {
		responseError(w, http.StatusNotFound, "item not found")
		return
	}
	response(w, http.StatusOK, rawData)
}

func (cs *CacheServer) removeHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len(cs.indexPath):]
	if key == "" {
		responseError(w, http.StatusBadRequest, "empty key passed")
		return
	}

	cs.c.Remove(key)
	w.WriteHeader(http.StatusOK)
	return

}

func (cs *CacheServer) statHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// do stuff to send back stat
			stats := map[string]interface{}{
				"size":    cs.c.Size(),
				"objects": cs.c.Len(),
			}
			responseJSON(w, http.StatusOK, stats)
			return
		default:
			w.WriteHeader(http.StatusNotImplemented)
			return
		}
	}
	return http.HandlerFunc(fn)
}

func responseError(w http.ResponseWriter, code int, message string) {
	response(w, code, []byte(message+"\n"))
}

func response(w http.ResponseWriter, code int, payload []byte) {
	w.WriteHeader(code)
	w.Write(payload)
}

func responseJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
