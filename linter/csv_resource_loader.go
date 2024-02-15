package linter

import (
	"github.com/stelligent/config-lint/assertion"
	"path/filepath"
        "fmt"
)

// CSVResourceLoader loads a list of Resource objects based on the list of ResourceConfig objects
type CSVResourceLoader struct {
	Columns []assertion.ColumnConfig
}

func extractCSVResourceID(expression string, properties interface{}) string {
	resourceID := "None"
	result, err := assertion.SearchData(expression, properties)
	if err == nil {
		resourceID, _ = result.(string)
	}
	return resourceID
}

// Load converts a text file into a collection of Resource objects
func (l CSVResourceLoader) Load(filename string) (FileResources, error) {
	loaded := FileResources{
		Resources: make([]assertion.Resource, 0),
	}
	csvRows, err := loadCSV(filename)
	if err != nil {
		return loaded, err
	}
	for rowNumber, row := range csvRows {
		properties := map[string]interface{}{}
		properties["__file__"] = filename
		properties["__dir__"] = filepath.Dir(filename)
		for columnNumber, columnConfig := range l.Columns {
			properties[columnConfig.Name] = row[columnNumber]
		}
		resource := assertion.Resource{
			ID:         fmt.Sprint(rowNumber),
			Type:       "row",
			Properties: properties,
			Filename:   filename,
		}
		loaded.Resources = append(loaded.Resources, resource)
	}
	return loaded, nil
}

// PostLoad does no additional processing fro a CSVResourceLoader
func (l CSVResourceLoader) PostLoad(r FileResources) ([]assertion.Resource, error) {
	return r.Resources, nil
}
