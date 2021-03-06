package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"searchQuery/elasticsearch"
	"searchQuery/log"
	"strings"

)

type healthMessage struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func HealthCheckHandlerCreator() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var (
			healthIssue string
			err         error
		)

		// assume all well
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		body := []byte("{\"status\":\"OK\"}") // quicker than json.Marshal(healthMessage{...})

		// test elastic access
		res, err := elasticsearch.GetStatus()
		if err != nil {
			healthIssue = err.Error()
		} else if !strings.Contains(string(res), " green ") {
			healthIssue = string(res)
		}

		// when there's a healthIssue, change headers and content
		if healthIssue != "" {
			w.WriteHeader(http.StatusInternalServerError)
			if body, err = json.Marshal(healthMessage{
				Status: "error",
				Error:  healthIssue,
			}); err != nil {
				log.Error(err, nil)
				panic(err)
			}
		}

		// return json
		fmt.Fprintf(w, string(body))
	}
}
