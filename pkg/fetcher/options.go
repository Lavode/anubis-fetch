package fetcher

import "time"

type Options struct {
	URL       string
	Timeout   time.Duration
	UA        string
	Text      bool
	Browser   bool // skip the solver, go straight to the browser
	NoBrowser bool // never use the browser; exit escalateExit if solving fails
	NoCache   bool // don't read/write the persistent cookie jar
}
