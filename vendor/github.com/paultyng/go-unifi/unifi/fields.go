package unifi

// This will generate the *.generated.go files in this package for the specified
// Unifi controller version.
//go:generate go run ../fields/ -version-base-dir=../fields/ -latest
