package openapi

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"

	"github.com/bmatcuk/doublestar/v4"
)

// MergeSpecs merges spec.Swagger specs into a single OpenAPI spec.
func MergeSpecs(title string, specs ...*spec.Swagger) ([]byte, error) {
	if len(specs) == 0 {
		return nil, fmt.Errorf("no specs provided")
	}

	specFinal := specs[0]
	analysis.Mixin(specFinal, specs...)
	specFinal.Info.Title = title

	return json.Marshal(specFinal)
}

// MergeDocuments merges loads.Documents into a single OpenAPI spec.
func MergeDocuments(title string, documents ...*loads.Document) ([]byte, error) {
	var specs []*spec.Swagger
	for _, document := range documents {
		specs = append(specs, document.Spec())
	}
	return MergeSpecs(title, specs...)
}

// MergeFiles merges OpenAPI files into a single OpenAPI spec.
func MergeFiles(title string, files ...string) ([]byte, error) {
	var specs []*loads.Document
	for _, file := range files {
		spec, err := loads.Spec(file)
		if err != nil {
			return nil, fmt.Errorf("error loading file: %w", err)
		}
		specs = append(specs, spec)
	}
	return MergeDocuments(title, specs...)
}

// MergeGlob merges OpenAPI files matching a glob pattern into a single OpenAPI
// spec.
func MergeGlob(title string, glob string) ([]byte, error) {
	matches, err := doublestar.Glob(os.DirFS("./"), glob)
	if err != nil {
		return nil, fmt.Errorf("error getting matches from glob: %w", err)
	}

	return MergeFiles(title, matches...)
}
