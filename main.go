package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gorilla/mux"
	"github.com/jaberchez/tom-frontend/pkg/backend"
)

var (
	backendService string
	podNamespace   string
)

func homePage(w http.ResponseWriter, r *http.Request) {
	backends, err := backend.Get(backendService, podNamespace)

	if err != nil {
		log.Println(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "some internal error ocurred")

		return
	}

	trunc := 80
	//data := make(map[string]string)
	var result string

	for _, item := range backends {
		//data := make(map[string]string)

		envs, err := backend.GetEnvVars(item.IP, item.Port)

		if err != nil {
			log.Println(err.Error())
			fmt.Fprintf(w, "<h1>Unable to get environment variables from backend %s (%s)</h1>", item.IP, item.Name)

			continue
		}

		if len(envs) == 0 {
			fmt.Fprintf(w, "<h1>Environment variables not found in backend %s (%s)</h1>", item.IP, item.Name)
			continue
		}

		keys := make([]string, 0, len(envs))

		for k := range envs {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		result = fmt.Sprintf(`
<html>
<head>
<style>
table, th, td {
  border: 1px solid black;
  border-collapse: collapse;
}

th, td {
	padding: 10px;
 }

tr:nth-child(even) {background-color: #f2f2f2;}
</style>
</head>
<body>

<h1>Environment Variables from backend %s (%s)</h1>
<table>
<tr>
<th>Name</th>
<th>Value</th>
</tr>
  `, item.IP, item.Name)

		for _, k := range keys {
			versionVariable := false

			result += "<tr>"

			if k == "APP_VERSION" {
				versionVariable = true
			}

			result += "<td>"

			if versionVariable {
				result += "<b>" + k + "</b>"
			} else {
				result += k
			}

			result += "</td>"

			if len(envs[k]) >= trunc {
				// Truncate the value
				val := envs[k]

				if versionVariable {
					result += "<td><b>" + val[:trunc] + "...</b></td>"
				} else {
					result += "<td>" + val[:trunc] + "...</td>"
				}
			} else {
				if versionVariable {
					result += "<td><b>" + envs[k] + "</b></td>"
				} else {
					result += "<td>" + envs[k] + "</td>"
				}
			}

			result += "</tr>"
		}

		result += "</table></html>"

		fmt.Fprintf(w, result)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func main() {
	backendService = os.Getenv("BACKEND_SERVICE")
	podNamespace = os.Getenv("POD_NAMESPACE")
	listenPort := os.Getenv("PORT")

	if len(backendService) == 0 {
		log.Fatal("unable to find the enviroment variable BACKEND_SERVICE")
	}

	if len(podNamespace) == 0 {
		log.Fatal("unable to find the enviroment variable POD_NAMESPACE")
	}

	if len(listenPort) == 0 {
		listenPort = "8080"
	}

	r := mux.NewRouter()

	r.HandleFunc("/", homePage).Methods("GET") // Only GETT allowed
	r.HandleFunc("/startup", healthCheck)
	r.HandleFunc("/liveness", healthCheck)
	r.HandleFunc("/readiness", healthCheck)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%s", listenPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server listening or port %s\n", listenPort)
	log.Fatal(srv.ListenAndServe())
}
