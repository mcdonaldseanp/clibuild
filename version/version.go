// EXAMPLE:
//
// This file is an example of a version package for a CLI tool built
// with clibuild. The clibuild.go example binary uses this version
// package to define the version reported with clibuild --version.
//
// The clibuild library provides a convenience command "[bin name] update version"
// that updates the version in this file to either 1. the provided version or 2.
// automatically bump the Z release.
package version

// When first starting out with a new project, this version should be
// v0.0.0
const VERSION string = "v0.0.3"
