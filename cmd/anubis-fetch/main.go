// Command anubis-fetch fetches a URL from behind a JavaScript bot-wall.
//
// It handles two very different walls with one tool:
//
//   - Anubis (github.com/TecharoHQ/anubis) — proof-of-work challenges.
//     The legacy SHA-256 solver and WASM runner both execute in-process.
//   - Cloudflare-style passive fingerprinting (TLS/JA3 + HTTP2) — cleared by
//     impersonating a real Chrome at the transport layer via curl-impersonate
//     (through the req library).
//
// Strategy, cheapest first:
//
//  1. Reuse a persisted cookie (a prior Anubis auth token) — like a browser,
//     a revisit isn't re-challenged.
//  2. Solve the proof-of-work in-process over the impersonating HTTP client.
//  3. Fall back to a real headless Chromium (chromedp) for anything the fast
//     path can't do: the preact/metarefresh challenge methods, an unknown or
//     future method, a difficulty too high to brute-force, a rejected solution,
//     or a Cloudflare *active* JS challenge.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/fzakaria/anubis-fetch/pkg/fetcher"
	"jaytaylor.com/html2text"
)

const (
	// EscalateExit is returned by --no-browser when the in-process solve fails.
	escalateExit = 3
)

func usage() {
	fmt.Fprint(os.Stderr, `anubis-fetch — fetch a URL from behind Anubis / Cloudflare bot-walls

Usage:
  anubis-fetch [flags] URL

Flags:
  --text            render readable plain text instead of HTML
  --timeout MS      per-step timeout in milliseconds (default 30000)
  --ua STRING       override the User-Agent
  --browser         skip the in-process solver; use the headless browser
  --no-browser      never use the browser; exit 3 if the solve can't apply
  --no-cache        don't read or write the persistent cookie jar
`)
}

func main() {
	os.Exit(run())
}

func run() int {
	o := parseFlags()

	if !o.Browser {
		html, escalate := fetcher.FetchViaHTTP(o)
		if !escalate {
			output(html, o)
			return 0
		}
		if o.NoBrowser {
			fmt.Fprintln(os.Stderr, "anubis-fetch: cannot solve in-process and --no-browser is set")
			return escalateExit
		}
	}

	html, err := fetcher.FetchViaBrowser(o)
	if err != nil {
		fmt.Fprintf(os.Stderr, "anubis-fetch: browser error: %v\n", err)
		if html == "" {
			return 1
		}
	}
	output(html, o)
	return 0
}

func parseFlags() fetcher.Options {
	var o fetcher.Options
	var timeoutMs int
	flag.IntVar(&timeoutMs, "timeout", 30000, "per-step timeout in milliseconds")
	flag.StringVar(&o.UA, "ua", "", "override User-Agent")
	flag.BoolVar(&o.Text, "text", false, "render readable plain text instead of HTML")
	flag.BoolVar(&o.Browser, "browser", false, "skip the solver; use the browser")
	flag.BoolVar(&o.NoBrowser, "no-browser", false, "never use the browser; exit 3 on failure")
	flag.BoolVar(&o.NoCache, "no-cache", false, "don't read or write the cookie jar")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		os.Exit(2)
	}
	o.URL = flag.Arg(0)
	o.Timeout = time.Duration(timeoutMs) * time.Millisecond
	return o
}

func output(html string, o fetcher.Options) {
	if o.Text {
		if text, err := html2text.FromString(html, html2text.Options{PrettyTables: true}); err == nil {
			fmt.Print(text)
			return
		}
		// Fall through to raw HTML if rendering fails.
	}
	fmt.Print(html)
}
