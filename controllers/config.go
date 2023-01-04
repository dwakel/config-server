package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/config-server/settings"
	"github.com/gorilla/mux"
	"sigs.k8s.io/yaml"
)

type Config struct {
	l *log.Logger
}

func NewConfig(l *log.Logger) *Config {
	return &Config{l}
}

func (this *Config) ServeConfig(rw http.ResponseWriter, r *http.Request) {
	//handle request
	vars := mux.Vars(r)

	path, _ := vars["path"]
	branch, _ := vars["branch"]

	//Fetch data
	client := http.Client{
		Timeout: 15 * time.Second,
	}
	owner := os.Getenv("OWNER")
	repo := os.Getenv("REPO_NAME")
	url := fmt.Sprintf("https://api.github.com/repos/%v/%v/contents/%v.yaml?ref=%v", owner, repo, path, branch)
	fmt.Println(url)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		this.l.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Failed to construct request")
		return
		//todo: Request failed. Handle it
	}
	basicAuthToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("Basic %v:%v", owner, os.Getenv("AUTH_TOKEN"))))
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Add("Authorization", fmt.Sprintf("Basic %v", basicAuthToken))
	request.Header.Set("X-GitHub-Api-Version", settings.GitHubAPIVersion)

	resp, err := client.Do(request)
	if err != nil {
		this.l.Println("Error")
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Server error")
		return
	}
	if resp.StatusCode > 201 {
		this.l.Println("Error")
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Server error")
		return
	}

	var results map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		this.l.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, err.Error())
		return
	}

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

	var source map[string]interface{}
	err = json.Unmarshal(jsonData, &source)
	if err != nil {
		this.l.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, err.Error())
		return
	}
	destination := make(map[string]interface{})
	for k, v := range source {
		if reflect.TypeOf(v) != reflect.TypeOf(map[string]interface{}{}) {
			destination[k] = v
			continue
		}
		for key, val := range v.(map[string]interface{}) {
			destination[k+"."+key] = val
		}
	}

	configResults, err := json.Marshal(destination)
	if err != nil {
		this.l.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, err.Error())
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(configResults))
	return

}
