package main

import (
	"os"

	"github.com/gophergala2016/supbot/lib/gitter"
	"github.com/gophergala2016/supbot/lib/slack"
)

func main() {
	slackToken := os.Getenv("SLACK_TOKEN")
	gitterToken := os.Getenv("GITTER_TOKEN")
	if slackToken == "" && gitterToken == "" {
		panic("Slack or Gitter token must be set")
	}

	if slackToken != "" {
		s := slack.NewClient(slackToken)
		s.Start()
	} else if gitterToken != "" {
		g, err := gitter.NewGitter(gitterToken)
		if err != nil {
			panic(err)
		}

		done := make(chan bool)
		g.Start(done)
		<-done // will wait forever
	}
}
