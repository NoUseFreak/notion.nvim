package cli

import (
	"context"

	"github.com/jomei/notionapi"
)

type CLI struct {
	client *notionapi.Client
	ctx    context.Context
}

func New(client *notionapi.Client, ctx context.Context) *CLI {
	return &CLI{
		client: client,
		ctx:    ctx,
	}
}
