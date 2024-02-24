package cli

import (
	"log"

	"github.com/jomei/notionapi"
)

func getPropByType(props map[string]notionapi.Property, propType notionapi.PropertyType) notionapi.Property {
	for key, prop := range props {
		if prop.GetType() == propType {
			return props[key]
		}
	}

	return nil
}

func getPropString(prop notionapi.Property) string {
	switch prop.GetType() {
	case notionapi.PropertyTypeTitle:
		for _, text := range prop.(*notionapi.TitleProperty).Title {
			return text.PlainText
		}
	case notionapi.PropertyTypeUniqueID:
		return prop.(*notionapi.UniqueIDProperty).UniqueID.String()
	default:
		log.Fatal("Unknown property type")
	}

	return ""
}
