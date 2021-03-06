package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/sontags/env"
	"github.com/unprofession-al/gerty/api"
	"github.com/unprofession-al/gerty/config"
	"github.com/unprofession-al/gerty/entities"
	mw "github.com/unprofession-al/gerty/middleware"
	"github.com/unprofession-al/gerty/store"
	"github.com/unprofession-al/gerty/transformers"
	// _ "github.com/unprofession-al/gerty/store/memstore"
	_ "github.com/unprofession-al/gerty/store/sqlitestore"
)

var cfg config.Configuration

func init() {
	env.Var(&cfg.Port, "PORT", "8008", "Port to bind to")
	env.Var(&cfg.Address, "ADDR", "0.0.0.0", "Address to bind to")
	env.Var(&cfg.Store, "STORE", "/tmp/gerty.sqlite3", "Store configuration string")
	env.Var(&cfg.NodeVarsProviders, "NODEVARS_PROVIDERS", "[]", "JSON string to configure nodevars providers")
	env.Var(&cfg.JenkinsFileName, "JEN_FILE_NAME", "inventory.json", "Name of the backup file")
	env.Var(&cfg.JenkinsToken, "JEN_TOKEN", "token", "Jenkins access token")
	env.Var(&cfg.JenkinsJobName, "JEN_JOB_NAME", "inventory.archive", "Jenkins job name")
	env.Var(&cfg.JenkinsBaseUrl, "JEN_BASE_URL", "NONE", "Jenkins base URL or 'NONE' in order to disable Jenkins WebHook")
}

func main() {
	env.Parse("GERTY", false)

	s, err := store.New("sqlitestore", cfg.Store)
	if err != nil {
		panic(err)
	}

	ri := entities.NewRoleInteractor(s.Roles)
	ni := entities.NewNodeInteractor(s.Nodes, ri)

	r := mux.NewRouter().StrictSlash(true)

	err = api.Configure(cfg)
	if err != nil {
		panic(err)
	}
	api.Inject(ni, ri)
	a := r.PathPrefix("/api/").Subrouter()
	api.PopulateRouter(a)

	transformers.Inject(ni, ri)
	t := r.PathPrefix("/transformers/").Subrouter()
	transformers.PopulateRouter(t)

	wh := mw.WebHook{
		FileName: cfg.JenkinsFileName,
		Token:    cfg.JenkinsToken,
		JobName:  cfg.JenkinsJobName,
		BaseUrl:  cfg.JenkinsBaseUrl,
	}

	chain := alice.New(
		mw.RecoverPanic,
		mw.CorsHeaders,
		mw.UserContext,
		wh.Create,
	).Then(r)

	log.Fatal(http.ListenAndServe(cfg.Address+":"+cfg.Port, chain))
}
