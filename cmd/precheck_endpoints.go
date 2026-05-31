//ff:func feature=scan type=engine control=iteration dimension=1
//ff:what Pre-checks existing hurl files for all endpoints and marks each session entry's initial verdict
package cmd

import (
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/session"
)

func precheckEndpoints(sess *session.Session, cfg *config.Config) {
	if cfg == nil {
		return
	}
	for i := range sess.Entries {
		precheckEntry(sess, &sess.Entries[i], cfg)
	}
}
