package cli

import (
	"context"

	"github.com/jomei/notionapi"
)

type CLI struct {
	client *notionapi.Client
	ctx    context.Context
	utils  *utils
}

func New(client *notionapi.Client, ctx context.Context) *CLI {
	return &CLI{
		client: client,
		ctx:    ctx,
		utils: &utils{
			ctx:    ctx,
			client: client,
		},
	}
}
