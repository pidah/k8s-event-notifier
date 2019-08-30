package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/bluele/slack"
	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
)

func main() {

	log.Print("Starting k8s-event-notifier...")

	var hookURL = os.Getenv("SLACK_API_URL")
	var slackText = os.Getenv("SLACK_TEXT")
	var eventFilter = strings.Split(os.Getenv("EVENT_FILTER"), ",")
	reasonList := []string{}
	reasonList = append(reasonList, eventFilter...)

	client, err := k8s.NewInClusterClient()
	if err != nil {
		log.Printf("Client ERROR: %s\n", err.Error())
	}

	k8s.Register("", "v1", "events", false, &corev1.Event{})
	var events corev1.Event

	for {

		watcher, err := client.Watch(context.Background(), "", &events)

		if err != nil {
			log.Printf("Watch ERROR: %s\n", err.Error())
		}

		defer watcher.Close()

	EventLoop:

		for {
			e := new(corev1.Event)
			_, err := watcher.Next(e)
			if err != nil {
				log.Printf("Event ERROR: %s\n", err.Error())
				watcher.Close()
				break EventLoop
			}
			if ok := stringInSlice(*e.Reason, reasonList); ok {
				log.Print(*e.InvolvedObject.Name, *e.Message)
				text := *e.InvolvedObject.Name + " in " + *e.InvolvedObject.Namespace + " namespace"
				hook := slack.NewWebHook(hookURL)
				err := hook.PostMessage(&slack.WebHookPostPayload{
					Text: slackText,
					Attachments: []*slack.Attachment{
						{Title: *e.Message, Text: text, Color: "#9741f4", Pretext: *e.Reason},
					},
				})
				if err != nil {
					log.Printf("Slack Post ERROR: %s\n", err.Error())
				}
			}
		}
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
