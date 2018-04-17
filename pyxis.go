package pyxis

import (
	"context"
	"encoding/json"
	"github.com/castaneai/asaka"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

const (
	baseURL           = "https://www.pixiv.net"
	sessionCookieName = "PHPSESSID"
)

type Client struct {
	hc        *http.Client
	sessionID string
}

func NewClient(hc *http.Client, sessionID string) (*Client, error) {
	return &Client{
		hc:        hc,
		sessionID: sessionID,
	}, nil
}

type Notification struct {
	Id         int       `json:"id" datastore:"id"`
	Content    string    `json:"content" datastore:"content,noindex"`
	NotifiedAt time.Time `json:"notifiedAt" datastore:"notifiedAt"`
	LinkURL    string    `json:"linkUrl" datastore:"linkUrl,noindex"`
	IconURL    string    `json:"iconUrl" datastore:"iconUrl,noindex"`
}

type responseJSONBody struct {
	Items []*Notification `json:"items"`
}

type responseJSON struct {
	Error   bool              `json:"error"`
	Message string            `json:"message"`
	Body    *responseJSONBody `json:"body"`
}

func createSessionCookie(session string) *http.Cookie {
	expires := time.Now().AddDate(1, 0, 0)
	return &http.Cookie{Name: sessionCookieName, Value: session, Expires: expires, HttpOnly: true}
}

func (c *Client) GetUsername(ctx context.Context) (string, error) {
	opts := &asaka.ClientOption{
		Cookies: map[string]http.Cookie{sessionCookieName: *createSessionCookie(c.sessionID)},
	}
	ac, err := asaka.NewClient(c.hc, opts)
	if err != nil {
		return "", err
	}
	doc, err := ac.GetDoc(ctx, baseURL)
	if err != nil {
		return "", err
	}
	return doc.Find(".user-name").Text(), nil
}

func (c *Client) GetNotifications(ctx context.Context) ([]*Notification, error) {
	req, err := http.NewRequest("GET", baseURL+"/ajax/notification", nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(createSessionCookie(c.sessionID))
	res, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var rj responseJSON
	if err := json.NewDecoder(res.Body).Decode(&rj); err != nil {
		return nil, err
	}
	if rj.Error {
		return nil, errors.New(rj.Message)
	}

	var result []*Notification
	// reverse slice
	for i := len(rj.Body.Items) - 1; i >= 0; i-- {
		result = append(result, rj.Body.Items[i])
	}
	return result, nil
}
