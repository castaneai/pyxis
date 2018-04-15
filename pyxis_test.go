package pyxis

import (
	"testing"
	"context"
	"net/http"
	"os"
)

func TestGetNotifications(t *testing.T) {
	sessionID := os.Getenv("PYXIS_SESSION")
	if sessionID == "" {
		t.Fatalf("env: PYXIS_SESSION not set")
	}
	c := &Client{hc: &http.Client{}, sessionID: sessionID}
	ctx := context.Background()
	ns, err := c.GetNotifications(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, n := range ns {
		t.Logf("%+v", n)
	}
}
