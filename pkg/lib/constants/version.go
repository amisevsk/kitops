package constants

import "runtime"

// Default build-time variables
// These values are overridden via ldflags
var (
	Version   = "unknown"
	GitCommit = "unknown"
	BuildTime = "unknown"
	GoVersion = runtime.Version()
)
