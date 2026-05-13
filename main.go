/*
Copyright © 2026 Shahid Amin aminShahid5515@gmail.com
*/
package main

import "gen/cmd"

// Version is set by goreleaser during build
var Version = "dev"

func main() {
	cmd.Execute()
}
