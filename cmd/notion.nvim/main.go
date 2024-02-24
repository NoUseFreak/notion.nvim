package main

import "github.com/nousefreak/notion.nvim/internal/app/cmd"

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	cmd.Version = Version
	cmd.Commit = Commit
	cmd.Date = Date
	cmd.Execute()
}
