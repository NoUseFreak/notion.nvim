package cli

import (
	"context"
	"fmt"
	"sync"

	"github.com/jomei/notionapi"
	"github.com/sirupsen/logrus"
)

type utils struct {
	ctx    context.Context
	client *notionapi.Client
}

func (u *utils) getPropByType(props map[string]notionapi.Property, propType notionapi.PropertyType) notionapi.Property {
	for key, prop := range props {
		if prop.GetType() == propType {
			return props[key]
		}
	}

	return nil
}

func (u *utils) getPropString(prop notionapi.Property) string {
	switch prop.GetType() {
	case notionapi.PropertyTypeTitle:
		for _, text := range prop.(*notionapi.TitleProperty).Title {
			return text.PlainText
		}
	case notionapi.PropertyTypeUniqueID:
		return prop.(*notionapi.UniqueIDProperty).UniqueID.String()
	case notionapi.PropertyTypeDate:
		date := prop.(*notionapi.DateProperty).Date
		if date != nil {
			return date.Start.String()
		}
		return ""
	default:
		logrus.Debugf("getPropString - Unknown property type: %s\n", prop.GetType())
	}

	return ""
}

func (u *utils) pageToIssue(page notionapi.Page) Issue {
	issue := Issue{
		ID:    u.getPropString(u.getPropByType(page.Properties, notionapi.PropertyTypeUniqueID)),
		Title: u.getPropString(u.getPropByType(page.Properties, notionapi.PropertyTypeTitle)),
		Assignees: (func() []string {
			assignees := []string{}
			for _, assignee := range u.getPropByType(page.Properties, notionapi.PropertyTypePeople).(*notionapi.PeopleProperty).People {
				assignees = append(assignees, assignee.Name)
			}
			return assignees
		})(),
		URL: page.URL,
	}

	for name, prop := range page.Properties {
		issueProp := IssueProperty{
			Name: name,
			Type: string(prop.GetType()),
		}
		switch prop.GetType() {
		case notionapi.PropertyTypeSelect:
			if value := u.getPropString(prop); value != "" {
				issueProp.Values = []string{value}
			}
		case notionapi.PropertyTypeRelation:
			for _, value := range prop.(*notionapi.RelationProperty).Relation {
				if title := u.getPageTitle(value.ID); title != nil {
					issueProp.Values = append(issueProp.Values, *title)
				}
			}
		case notionapi.PropertyTypeMultiSelect:
			for _, value := range prop.(*notionapi.MultiSelectProperty).MultiSelect {
				issueProp.Values = append(issueProp.Values, value.Name)
			}
		case notionapi.PropertyTypePeople:
			for _, value := range prop.(*notionapi.PeopleProperty).People {
				issueProp.Values = append(issueProp.Values, value.Name)
			}
		case notionapi.PropertyTypeDate:
			if data := u.getPropString(prop); data != "" {
				issueProp.Values = []string{data}
			}
		case notionapi.PropertyTypeTitle, notionapi.PropertyTypeUniqueID, notionapi.PropertyTypeStatus:
			// Do nothing
		default:
			logrus.Debugf("Unknown property type: %s", prop.GetType())
		}

		if len(issueProp.Values) > 0 {
			issue.Properties = append(issue.Properties, issueProp)
		}
	}
	return issue
}

var pageTitleCache = sync.Map{}

func (u *utils) getPageTitle(id notionapi.PageID) *string {
	if title, ok := pageTitleCache.Load(id); ok {
		return title.(*string)
	}

	page, err := u.client.Page.Get(u.ctx, id)
	if err != nil {
		logrus.Errorf("Error fetching page: %s", id)
		return nil
	}

	for _, prop := range page.Properties {
		if prop.GetType() == notionapi.PropertyTypeTitle {
			title := u.getPropString(prop)
			pageTitleCache.Store(id, title)
			return &title
		}
	}
	return nil
}

func (u *utils) markdownifyBlock(block notionapi.Block) string {
	switch block.GetType() {
	case notionapi.BlockTypeParagraph:
		return block.GetRichTextString()
	case notionapi.BlockTypeHeading1:
		return fmt.Sprintf("# %s", block.GetRichTextString())
	case notionapi.BlockTypeHeading2:
		return fmt.Sprintf("## %s", block.GetRichTextString())
	case notionapi.BlockTypeHeading3:
		return fmt.Sprintf("### %s", block.GetRichTextString())
	case notionapi.BlockTypeBulletedListItem:
		return fmt.Sprintf("* %s", block.GetRichTextString())
	case notionapi.BlockTypeNumberedListItem:
		return fmt.Sprintf("1. %s", block.GetRichTextString())
	case notionapi.BlockTypeToDo:
		return fmt.Sprintf("- [ ] %s", block.GetRichTextString())
	case notionapi.BlockTypeToggle:
		return fmt.Sprintf("::: details\n%s\n:::", block.GetRichTextString())
	case notionapi.BlockTypeChildPage:
		return fmt.Sprintf("[%s](%s)", block.GetRichTextString(), block.GetRichTextString())
	case notionapi.BlockTypeUnsupported:
		return ""
	default:
		return ""
	}
}
