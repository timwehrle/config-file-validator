package validator

import (
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"

	"github.com/lestrrat-go/libxml2"
)

type XMLValidator struct{}

// Validate implements the Validator interface by attempting to
// unmarshall a byte array of xml
func (XMLValidator) Validate(b []byte) (bool, error) {
	var output any
	err := xml.Unmarshal(b, &output)
	if err != nil {
		return false, err
	}

	doc, err := libxml2.ParseString(string(b))
	if err != nil {
		return false, err
	}
	defer doc.Free()

	schemaLocation, err := extractSchemaLocation(string(b))
	if err != nil {
		return false, err
	}

	fmt.Println(schemaLocation)

	return true, nil
}

func extractSchemaLocation(xmlContent string) (string, error) {
	schemaLocationRegex := regexp.MustCompile(`xsi:schemaLocation\s*=\s*"([^"]+)"`)
	noNamespaceSchemaLocationRegex := regexp.MustCompile(`xsi:noNamespaceSchemaLocation\s*=\s*"([^"]+)"`)

	matches := schemaLocationRegex.FindStringSubmatch(xmlContent)
	if len(matches) > 1 {
		return matches[1], nil
	}

	noNsMatches := noNamespaceSchemaLocationRegex.FindStringSubmatch(xmlContent)
	if len(noNsMatches) > 1 {
		return noNsMatches[1], nil
	}

	return "", errors.New("no schema location found")
}