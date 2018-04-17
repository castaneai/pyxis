package main

import (
	"net/http"

	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/castaneai/pyxis"
	"github.com/grokify/html-strip-tags-go"
	"github.com/pkg/errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	http.HandleFunc("/", handle)
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	hc := urlfetch.Client(ctx)
	cli, err := pyxis.NewClient(hc, os.Getenv("PYXIS_SESSION"))
	if err != nil {
		log.Errorf(ctx, "%+v", err)
		return
	}

	ns, err := cli.GetNotifications(ctx)
	if err != nil {
		log.Errorf(ctx, "%+v", err)
		return
	}
	for _, n := range ns {
		fmt.Fprintf(w, "<div><img src=\"%s\">%s</div>", n.IconURL, n.Content)
	}

	last, err := getLastNotification(ctx)
	if err != nil {
		log.Errorf(ctx, "%+v", err)
		return
	}
	fns := filterNewNotifications(ns, last)
	if len(fns) > 0 {
		if err := saveNewNotifications(ctx, fns); err != nil {
			log.Errorf(ctx, "%+v", err)
			return
		}
		if err := postNewNotificationsToSlack(ctx, fns); err != nil {
			log.Errorf(ctx, "%+v", err)
			return
		}
		fmt.Fprintf(w, "<h2>%d notifications saved.<h2>", len(fns))
	}
}

func postNewNotificationsToSlack(ctx context.Context, ns []*pyxis.Notification) error {
	slackWebhookURL := os.Getenv("PYXIS_SLACK_WEBHOOK_URL")
	if slackWebhookURL == "" {
		return errors.New("env: PYXIS_SLACK_WEBHOOK_URL not set")
	}
	for _, n := range ns {
		if err := postNotificationToSlack(ctx, slackWebhookURL, n); err != nil {
			return err
		}
	}
	return nil
}

func filterNewNotifications(ns []*pyxis.Notification, lastNotification *pyxis.Notification) []*pyxis.Notification {
	var fns []*pyxis.Notification
	for _, n := range ns {
		if isNewNotification(n, lastNotification) {
			fns = append(fns, n)
		}
	}
	return fns
}

func isNewNotification(n *pyxis.Notification, last *pyxis.Notification) bool {
	if last == nil {
		return true
	}
	return n.Id != last.Id && !n.NotifiedAt.Before(last.NotifiedAt)
}

const (
	DatastoreKind = "Notifications"
)

func getLastNotification(ctx context.Context) (*pyxis.Notification, error) {
	q := datastore.NewQuery(DatastoreKind).Order("-notifiedAt").Limit(1)

	var lastNotifications []*pyxis.Notification
	if _, err := q.GetAll(ctx, &lastNotifications); err != nil {
		return nil, err
	}
	if len(lastNotifications) < 1 {
		return nil, nil
	}
	return lastNotifications[0], nil
}

func saveNewNotifications(ctx context.Context, notifications []*pyxis.Notification) error {
	var keys []*datastore.Key
	for range notifications {
		keys = append(keys, datastore.NewIncompleteKey(ctx, DatastoreKind, nil))
	}
	return datastore.RunInTransaction(ctx, func(tc context.Context) error {
		if _, err := datastore.PutMulti(ctx, keys, notifications); err != nil {
			return err
		}
		return nil
	}, nil)
}

type SlackMessage struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

func postSlackMessage(ctx context.Context, slackWebHookURL string, message *SlackMessage) error {
	hc := urlfetch.Client(ctx)
	buf, err := json.Marshal(message)
	if err != nil {
		return err
	}
	resp, err := hc.Post(slackWebHookURL, "application/json", bytes.NewReader(buf))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func postNotificationToSlack(ctx context.Context, slackWebhookURL string, notification *pyxis.Notification) error {
	text := strip.StripTags(notification.Content)
	mes := &SlackMessage{
		Text:    text,
		IconURL: notification.IconURL,
	}
	err := postSlackMessage(ctx, slackWebhookURL, mes)
	if err != nil {
		return err
	}
	return nil
}
