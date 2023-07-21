package pkg

import (
    "strings"
    "regexp"
	"github.com/google/uuid"
	"text/template"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
var removeSpace   = regexp.MustCompile("( )")

func ToSnakeCase(str string) string {
    snake := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
    snake  = matchAllCap.ReplaceAllString(snake, "${1}-${2}")
    snake  = removeSpace.ReplaceAllString(snake, "${2}")
    return strings.ToLower(snake)
}

// Define a custom function to capitalize the first letter of a string.
func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}

func appId(len int) string {
	id := uuid.New()
    val := id.String()
	return strings.ReplaceAll(strings.ToUpper(val), "-", "")[0:len]
}


func NewTemplate() (*template.Template) {
	// Create a new template and register the custom function.
	tmpl := template.New("example").Funcs(template.FuncMap{
		"capitalize": capitalizeFirstLetter,
		"appId": appId,
		"snakecase": ToSnakeCase,
		"tolower": strings.ToLower,
	})
	return tmpl
}