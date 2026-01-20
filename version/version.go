package version

// Version is the current version of sheeit.
// This can be overridden at build time using:
//
//	go build -ldflags "-X github.com/zomglings/sheeit/version.Version=1.0.0"
var Version = "0.0.1"
