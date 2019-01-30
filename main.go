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
	var eventFilter = strings.Split(os.Getenv("EVENT_FILTER"), ",")
	reasonList := []string{}
	reasonList = append(reasonList, eventFilter...)

	client, err := k8s.NewInClusterClient()
	if err != nil {
		log.Fatal(err)
	}

	k8s.Register("", "v1", "events", false, &corev1.Event{})
	var events corev1.Event
	watcher, err := client.Watch(context.Background(), "", &events)

	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()
	for {
		e := new(corev1.Event)
		_, err := watcher.Next(e)
		if err != nil {
			log.Fatal(err)
		}
		if ok := stringInSlice(*e.Reason, reasonList); ok {
			log.Print(*e.InvolvedObject.Name, *e.Message)
			hook := slack.NewWebHook(hookURL)
			err := hook.PostMessage(&slack.WebHookPostPayload{
				Text: "cluster autoscaler:",
				Attachments: []*slack.Attachment{
					{Title: *e.Message, Text: *e.InvolvedObject.Name, Color: "#9741f4", Pretext: *e.Reason},
				},
			})
			if err != nil {
				log.Print(err)
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
