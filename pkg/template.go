package pkg

import (
	"github.com/google/uuid"
	"regexp"
	"strings"
	"text/template"
)

// var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
// var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
// var removeSpace = regexp.MustCompile("( )")

// KebabCase converts a string like "hello world" to "hello-world"
func KebabCase(str string) string {
	// snake := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	// snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")
	// snake = removeSpace.ReplaceAllString(snake, "${2}")
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	result := re.ReplaceAllString(str, "-")
	return strings.ToLower(result)
}

// StartCase defines a custom function to capitalize the first
// letter of a string.
func StartCase(str string) string {
	if len(str) == 0 {
		return ""
	}

	return strings.ToUpper(str[0:1]) + strings.ToLower(str[1:])
}

// AppId generates a unique application ID
func AppId(length int) string {
	id := uuid.New()
	val := id.String()
	return strings.ReplaceAll(strings.ToUpper(val), "-", "")[0:length]
}

// NewTemplate is a factory to generate a template parser with the following
// helper methods embedded:
//   - capitalize
//   - appId
//   - kebabCase
//   - toLower
//
// which can then be used like so in a `tmpl` file:
//
//	name: {{ kebabCase .Service }}-dev
func NewTemplate() *template.Template {
	tmpl := template.New("configGenerator").Funcs(template.FuncMap{
		"capitalize": StartCase,
		"appId":      AppId,
		"kebabCase":  KebabCase,
		"toLower":    strings.ToLower,
	})
	return tmpl
}
