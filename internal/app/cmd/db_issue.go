package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gomarkdown/markdown/parser"
	"github.com/jomei/notionapi"
	"github.com/nousefreak/notion.nvim/internal/app/cli"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	markdown "github.com/MichaelMure/go-term-markdown"
	md "github.com/gomarkdown/markdown"
)

func getDBIssueDetailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db-issue-detail",
		Short: "db-issue-detail is a command line tool for interacting with the Notion API",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			token := os.Getenv("NOTION_INTEGRATION_TOKEN")
			dbID, _ := cmd.Flags().GetString("db-id")
			issueID := args[0]

			cli := cli.New(notionapi.NewClient(notionapi.Token(token)), cmd.Context())
			issue, err := cli.GetIssue(
				dbID,
				issueID,
			)
			if err != nil {
				log.Fatal(err)
			}

			if renderContent, _ := cmd.Flags().GetBool("render-content"); renderContent {
				props := `
ID
: ` + issue.ID + `

Title
: ` + issue.Title + `

Assignees
: ` + strings.Join(issue.Assignees, "\n:  ") + `

URL
: ` + issue.URL + `

`

				meta := ""
				for _, prop := range issue.Properties {
					if len(prop.Values) != 0 {
						meta += fmt.Sprintf(": __%s__: %s\n", prop.Name, strings.Join(prop.Values, ", "))
					}
				}
				if meta != "" {
					props += "Properties\n" + meta
				}
				props += "\n\n---\n\n"

				content := props + strings.Join(issue.Content, "\n")
				width, _, err := term.GetSize(0)
				if err != nil {
					width = 80
				}

				exts := markdown.Extensions()
				exts ^= parser.Autolink
				p := parser.NewWithExtensions(exts)
				nodes := md.Parse([]byte(content), p)
				renderer := markdown.NewRenderer(width, 1)

				fmt.Println(string(md.Render(nodes, renderer)))
			} else {
				data, err := json.Marshal(issue)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(string(data))
			}
		},
	}

	cmd.Flags().String("db-id", "", "The ID of the database to query")
	cmd.Flags().Bool("render-content", false, "Render content blocks")
	cmd.MarkFlagRequired("db-id")

	return cmd
}

func getDBIssueCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db-issue",
		Short: "db-issue is a command line tool for interacting with the Notion API",
		Run: func(cmd *cobra.Command, args []string) {
			token := os.Getenv("NOTION_INTEGRATION_TOKEN")
			dbID, _ := cmd.Flags().GetString("db-id")
			owner, _ := cmd.Flags().GetString("owner")
			includeClosed, _ := cmd.Flags().GetBool("include-closed")

			input := cli.GetIssuesInput{
				Search:        strings.Join(args, " "),
				User:          owner,
				IncludeClosed: includeClosed,
			}

			cli := cli.New(notionapi.NewClient(notionapi.Token(token)), cmd.Context())
			issues, err := cli.GetIssues(
				dbID,
				input,
			)
			if err != nil {
				logrus.Fatal(err)
			}

			start := time.Now()
			data, err := json.Marshal(issues)
			if err != nil {
				logrus.Fatal(err)
			}
			fmt.Println(string(data))
			logrus.Debugf("Format took: %s", time.Since(start))
		},
	}

	cmd.Flags().String("db-id", "", "The ID of the database to query")
	cmd.Flags().String("owner", "", "Filter issues assigned to this user")
	cmd.Flags().Bool("include-closed", false, "Include closed issues")

	cmd.MarkFlagRequired("db-id")

	return cmd
}
