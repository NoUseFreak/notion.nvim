package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jomei/notionapi"
	"github.com/sirupsen/logrus"
)

func (c *CLI) getIssueDBMap(dbID string) (IssueDBSpec, error) {
	start := time.Now()
	loadCache()

	if cached, err := diskCache.Get(dbID); err == nil {
		var spec IssueDBSpec
		if err := json.Unmarshal([]byte(cached), &spec); err == nil {
			return spec, nil
		}
	}

	db, err := c.client.Database.Get(c.ctx, notionapi.DatabaseID(dbID))
	if err != nil {
		return IssueDBSpec{}, err
	}

	propMap := IssueDBSpec{}
	for name, prop := range db.Properties {
		if v, ok := prop.(*notionapi.UniqueIDPropertyConfig); ok {
			propMap.ID = name
			propMap.IDPrefix = v.UniqueID.Prefix
		}
		if _, ok := prop.(*notionapi.TitlePropertyConfig); ok {
			propMap.Title = name
		}
		if _, ok := prop.(*notionapi.PeoplePropertyConfig); ok {
			propMap.Assignees = name
		}
		if _, ok := prop.(*notionapi.StatusPropertyConfig); ok {
			propMap.Status = name
		}
	}

	if propMap.ID == "" || propMap.Title == "" {
		return IssueDBSpec{}, fmt.Errorf("Database does not contain required properties")
	}

	if specJSON, err := json.Marshal(propMap); err == nil {
		diskCache.Set(dbID, string(specJSON))
	}

	logrus.Debugf("getIssueDBMap took: %s", time.Since(start))

	return propMap, nil
}

func (c *CLI) GetIssue(dbID string, issueID string) (Issue, error) {
	propMap, err := c.getIssueDBMap(dbID)
	if err != nil {
		return Issue{}, err
	}

	re := regexp.MustCompile(fmt.Sprintf("^%s-([0-9]+)$", propMap.IDPrefix))
	if !re.MatchString(issueID) {
		return Issue{}, fmt.Errorf("Invalid issue ID")
	}

	number, err := strconv.Atoi(re.FindStringSubmatch(issueID)[1])
	if err != nil {
		return Issue{}, err
	}
	floatNumber := float64(number)

	response, err := c.client.Database.Query(c.ctx, notionapi.DatabaseID(dbID), &notionapi.DatabaseQueryRequest{
		Filter: notionapi.PropertyFilter{
			Property: propMap.ID,
			Number: &notionapi.NumberFilterCondition{
				Equals: &floatNumber,
			},
		},
		PageSize: 1,
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(response.Results) == 0 {
		return Issue{}, fmt.Errorf("Issue not found")
	}

	result := response.Results[0]
	issue := c.utils.pageToIssue(result)

	children, err := c.client.Block.GetChildren(c.ctx, notionapi.BlockID(result.ID), nil)
	if err != nil {
		log.Fatal(err)
	}
	content := []string{}
	for _, block := range children.Results {
		content = append(content, c.utils.markdownifyBlock(block))
	}
	issue.Content = content

	return issue, nil
}

type GetIssuesInput struct {
	Search        string
	User          string
	IncludeClosed bool
}

func (c *CLI) GetIssues(dbID string, input GetIssuesInput) ([]Issue, error) {
	propMap, err := c.getIssueDBMap(dbID)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	re := regexp.MustCompile(fmt.Sprintf("^%s-([0-9]+)$", propMap.IDPrefix))

	var filter notionapi.AndCompoundFilter
	var filterParts notionapi.AndCompoundFilter
	for _, part := range strings.Split(strings.TrimSpace(input.Search), " ") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		filterParts = append(filterParts, &notionapi.PropertyFilter{
			Property: propMap.Title,
			RichText: &notionapi.TextFilterCondition{
				Contains: part,
			},
		})

		if re.MatchString(part) {
			if number, err := strconv.Atoi(re.FindStringSubmatch(part)[1]); err == nil {
				floatNumber := float64(number)
				filter = append(filter, &notionapi.PropertyFilter{
					Property: propMap.ID,
					Number: &notionapi.NumberFilterCondition{
						Equals: &floatNumber,
					},
				})
			}
		}

		if number, err := strconv.Atoi(part); err == nil {
			floatNumber := float64(number)
			filter = append(filter, &notionapi.PropertyFilter{
				Property: propMap.ID,
				Number: &notionapi.NumberFilterCondition{
					Equals: &floatNumber,
				},
			})
		}
	}

	if input.User != "" {
		filter = append(filter, notionapi.PropertyFilter{
			Property: propMap.Assignees,
			People: &notionapi.PeopleFilterCondition{
				Contains: input.User,
			},
		})
	}

	if !input.IncludeClosed {
		filter = append(filter, &notionapi.PropertyFilter{
			Property: propMap.Status,
			Status: &notionapi.StatusFilterCondition{
				DoesNotEqual: "done",
			},
		})
	}

	if len(filterParts) > 0 {
		filter = append(filter, &filterParts)
	}

	logrus.Debugf("Filter constructed in: %s", time.Since(start))
	return c.doQuery(dbID, filter)
}

func (c *CLI) doQuery(dbID string, filter notionapi.AndCompoundFilter) ([]Issue, error) {
	start := time.Now()
	limit := 100
	var response *notionapi.DatabaseQueryResponse
	var err error
	if len(filter) == 0 {
		response, err = c.client.Database.Query(c.ctx, notionapi.DatabaseID(dbID), &notionapi.DatabaseQueryRequest{
			PageSize: limit,
		})
	} else {
		response, err = c.client.Database.Query(c.ctx, notionapi.DatabaseID(dbID), &notionapi.DatabaseQueryRequest{
			Filter:   filter,
			PageSize: limit,
		})
	}
	if err != nil {
		log.Fatal(err)
	}
	logrus.Debugf("Query took: %s", time.Since(start))

	start = time.Now()
	issues := []Issue{}

	channel := make(chan Issue)
	for _, result := range response.Results {
		go func(result notionapi.Page) {
			start := time.Now()
			issue := c.utils.pageToIssue(result)
			logrus.Debugf("Processing item %s took %s", issue.ID, time.Since(start))
			channel <- issue
		}(result)
	}

	for range response.Results {
		issues = append(issues, <-channel)
	}

	logrus.Debugf("Processing took: %s", time.Since(start))

	return issues, nil
}
