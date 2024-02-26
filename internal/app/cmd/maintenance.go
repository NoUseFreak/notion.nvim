package cmd

import (
	"fmt"
	"os"

	"github.com/nousefreak/notion.nvim/internal/app/cache"
	"github.com/nousefreak/notion.nvim/internal/app/cli"
	"github.com/spf13/cobra"
)

func getMaintenanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "maintenance",
		Run: func(cmd *cobra.Command, args []string) {
			path := cache.GetDBPath(cli.CACHE_KEY)
			if err := os.Remove(path); err != nil {
				panic(err)
			}

			fmt.Printf("Removed cache at: %s\n", path)
		},
	}

	return cmd
}
