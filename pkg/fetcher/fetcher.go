package fetcher

const (
	passChallengePath = "/.within.website/x/cmd/anubis/api/pass-challenge"
	// EscalateExit is returned by --no-browser when the in-process solve fails.
	EscalateExit = 3
	// maxDifficulty caps the in-process brute force. Go's SHA-256 is fast
	// (~16M hashes ≈ a few seconds at difficulty 6); above that a browser's
	// parallel workers win, so we escalate instead.
	maxDifficulty = 6
	defaultUA     = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36"
)

// legacyMethods use the native SHA-256 solver and its hex-digit difficulty cap.
// WASM methods use the server module and a time budget instead.
var legacyMethods = map[string]bool{"fast": true, "slow": true}
