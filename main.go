package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "probe" {
		client := http.Client{Timeout: 2 * time.Second}
		res, err := client.Get("http://127.0.0.1:8080/health")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			log.Fatal(res.Status)
		}
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /now", func(w http.ResponseWriter, r *http.Request) {
		zone := r.URL.Query().Get("tz")
		if zone == "" {
			zone = "UTC"
		}
		loc, err := time.LoadLocation(zone)
		if err != nil {
			http.Error(w, "time zone unavailable", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"tz": zone, "time": time.Now().In(loc).Format(time.RFC3339)})
	})
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("HEALTH_FAIL") == "1" {
			http.Error(w, "forced failure", 500)
			return
		}
		_, _ = w.Write([]byte("ok"))
	})
	server := http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(server.ListenAndServe())
}
