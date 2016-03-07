package ape

import (
	"fmt"
	"net/http"

	"github.com/groob/ape/api"
	"github.com/mholt/caddy/caddy/setup"
	"github.com/mholt/caddy/middleware"
)

type handler struct {
	Paths      []string
	Next       middleware.Handler
	apihandler http.Handler
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	// if the request path is any of the configured paths
	// write hello
	for _, p := range h.Paths {
		if middleware.Path(r.URL.Path).Matches(p) {
			h.apihandler.ServeHTTP(w, r)
			return 0, nil
		}
	}
	return h.Next.ServeHTTP(w, r)
}

// Setup creates a caddy middleware
func Setup(c *setup.Controller) (middleware.Middleware, error) {
	repoPath, err := parseRepo(c)
	if err != nil {
		return nil, err
	}
	paths := []string{"/api"}

	// Runs on Caddy startup, useful for services or other setups.
	c.Startup = append(c.Startup, func() error {
		fmt.Println("api middleware is initiated")
		return nil
	})

	// Runs on Caddy shutdown, useful for cleanups.
	c.Shutdown = append(c.Shutdown, func() error {
		fmt.Println("api middleware is cleaning up")
		return nil
	})

	return func(next middleware.Handler) middleware.Handler {

		h := &handler{
			Paths: paths,
			Next:  next,
		}
		repo := api.SimpleRepo(repoPath)
		server := api.NewServer(repo)
		h.apihandler = server
		return h
	}, nil
}

func parseRepo(c *setup.Controller) (string, error) {
	var repo string
	for c.Next() {
		args := c.RemainingArgs()
		switch len(args) {
		case 0:
			for c.NextBlock() {
				switch c.Val() {
				case "repo":
					repo = c.Val()
				}
			}
		case 1:
			repo = args[0]
		default:
			return repo, c.ArgErr()
		}
	}
	return repo, nil
}
