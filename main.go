package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/google/go-github/v32/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var client *github.Client
var token = envOrString("GH_PAT", "")

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Error("error reading request body")
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Error("could not parse webhook")
		return
	}

	switch e := event.(type) {
	case *github.PushEvent:
		if e.GetRef() == "refs/heads/main" {
			log.Info("Push on main!")
			ctx := context.Background()
			opts := github.DispatchRequestOptions{
				EventType:     "update-version",
				ClientPayload: nil,
			}
			_, _, err := client.Repositories.Dispatch(ctx, "camal-cakar-gcx", "gh-adapter-deployment-stub", opts)
			if err != nil {
				log.WithFields(log.Fields{
					"Error": err.Error(),
				}).Error("could not send dispatch event")
			}
		} else {
			log.Info("Push on unknown branch!")
		}
	default:
		log.Info("Unknown Event!")
	}
}

func envOrString(key, defaultValue string) (value string) {
	value, found := os.LookupEnv(key)
	if found {
		return
	}
	value = defaultValue
	return
}

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client = github.NewClient(tc)

	log.Info("Server started listening on :8080")
	http.HandleFunc("/webhook", handleWebhook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
