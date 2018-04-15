package main

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"github.com/castaneai/pyxis"
	"os"
	"fmt"
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
}

