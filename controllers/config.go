package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/config-server/settings"
	"github.com/gorilla/mux"
	"sigs.k8s.io/yaml"
)

type Config struct {
	l         *log.Logger
	owner     string
	repo      string
	authToken string
}

func NewConfig(l *log.Logger, owner, repo, authToken string) *Config {
	return &Config{l, owner, repo, authToken}
}

func (this *Config) ServeConfig(rw http.ResponseWriter, r *http.Request) {
	//handle request
	vars := mux.Vars(r)
	path, exist := vars["path"]
	if !exist {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(rw, "Not Found")
		return
	}
	branch, exist := vars["branch"]
	if !exist {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(rw, "Not Found")
		return
	}

	//Fetch data
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	url := fmt.Sprintf(settings.GitHubURL, this.owner, this.repo, path, branch)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		this.l.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Failed to construct request")
		return
		//todo: Request failed. Handle it
	}
	basicAuthToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("Basic %v:%v", this.owner, this.authToken)))
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Add("Authorization", fmt.Sprintf("Basic %v", basicAuthToken))
	request.Header.Set("X-GitHub-Api-Version", settings.GitHubAPIVersion)

	resp, err := client.Do(request)
	if err != nil {
		this.l.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Server error")
		return
	}
	if resp.StatusCode > http.StatusOK {
		this.l.Println("Request Failed")
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Server error")
		return
	}

	var results map[string]any
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		this.l.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, err.Error())
		return
	}

	//extracting github file from base64 string
	str := results["content"].(string)
	decText, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		this.l.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, err.Error())
		return
	}

	jsonData, err := yaml.YAMLToJSON(decText)
	if err != nil {
		this.l.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, err.Error())
		return
	}

	var source map[string]any
	err = json.Unmarshal(jsonData, &source)
	if err != nil {
		this.l.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, err.Error())
		return
	}

	//Flattern out JSON
	flatJson := flatternJson(source)

	configResults, err := json.Marshal(flatJson)
	if err != nil {
		this.l.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, err.Error())
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(configResults))
	return

}

func flatternJson(source map[string]any) map[string]any {
	//todo: Find a more elegant way to do this
	flatten := make(map[string]any)
	for k, v := range source {
		if reflect.TypeOf(v) != reflect.TypeOf(map[string]any{}) {
			flatten[k] = v
			continue
		}
		for key, val := range v.(map[string]any) {
			flatten[k+"."+key] = val
		}
	}
	return flatten
}
