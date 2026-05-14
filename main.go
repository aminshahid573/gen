/*
Copyright © 2026 Shahid Amin aminShahid5515@gmail.com
*/
package main

import (
    "runtime/debug"
    "github.com/aminshahid573/gen/cmd"
)

var version = "dev"

func main() {
    if version == "dev" {
        if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
            version = info.Main.Version
        }
    }
    cmd.Execute(version)
}
