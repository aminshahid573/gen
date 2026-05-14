/*
Copyright © 2026 Shahid Amin aminShahid5515@gmail.com
*/
package main

import "gen/cmd"

// version is set by goreleaser during build
var version = "dev"

func main() {
	cmd.Execute(version)
}
