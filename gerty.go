package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/sontags/env"
	"github.com/unprofession-al/gerty/api"
	"github.com/unprofession-al/gerty/entities"
	mw "github.com/unprofession-al/gerty/middleware"
	"github.com/unprofession-al/gerty/store"
	"github.com/unprofession-al/gerty/transformers"
	// _ "github.com/unprofession-al/gerty/store/memstore"
	_ "github.com/unprofession-al/gerty/store/sqlitestore"
)

type configuration struct {
	Port            string `json:"port"`
	Address         string `json:"address"`
	Store           string `json:"store"`
	JenkinsFileName string `json:"jenkins_file_name"`
	JenkinsToken    string `json:"-"`
	JenkinsJobName  string `json:"jenkins_job_name"`
	JenkinsBaseUrl  string `json:"jenkins_base_url"`
}

var config configuration

func init() {
	env.Var(&config.Port, "PORT", "8008", "Port to bind to")
	env.Var(&config.Address, "ADDR", "0.0.0.0", "Address to bind to")
	env.Var(&config.Store, "STORE", "/tmp/gerty.sqlite3", "Store configuration string")
	env.Var(&config.JenkinsFileName, "JEN_FILE_NAME", "inventory.json", "Name of the backup file")
	env.Var(&config.JenkinsToken, "JEN_TOKEN", "token", "Jenkins access token")
	env.Var(&config.JenkinsJobName, "JEN_JOB_NAME", "inventory.archive", "Jenkins job name")
	env.Var(&config.JenkinsBaseUrl, "JEN_BASE_URL", "NONE", "Jenkins base URL or 'NONE' in order to disable Jenkins WebHook")
}

func main() {
	env.Parse("GERTY", false)

	s, err := store.New("sqlitestore", config.Store)
	if err != nil {
		panic(err)
	}

	ri := entities.NewRoleInteractor(s.Roles)
	ni := entities.NewNodeInteractor(s.Nodes, ri)

	r := mux.NewRouter().StrictSlash(true)

	api.Inject(ni, ri)
	a := r.PathPrefix("/api/").Subrouter()
	api.PopulateRouter(a)

	transformers.Inject(ni, ri)
	t := r.PathPrefix("/transformers/").Subrouter()
	transformers.PopulateRouter(t)

	wh := mw.WebHook{
		FileName: config.JenkinsFileName,
		Token:    config.JenkinsToken,
		JobName:  config.JenkinsJobName,
		BaseUrl:  config.JenkinsBaseUrl,
	}

	chain := alice.New(
		mw.RecoverPanic,
		mw.CorsHeaders,
		mw.UserContext,
		wh.Create,
	).Then(r)

	log.Fatal(http.ListenAndServe(config.Address+":"+config.Port, chain))
}
