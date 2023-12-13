package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/charmbracelet/lipgloss"
)

// Styling for the version check UI
var (
	box = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).
		MarginTop(1).
		PaddingLeft(2).
		PaddingRight(2).
		PaddingTop(1).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("3"))
)

type ReleaseJson struct {
	TagName string `json:"tag_name"`
}

const releaseUrl = "https://api.github.com/repos/bcgov/gwa-cli/releases/latest"

// requestPublishedVersion pings the GitHub release API to
// retrieve the current version
func requestPublishedVersion() (string, error) {
	req, err := http.NewRequest("GET", releaseUrl, nil)
	if err != nil {
		return "", err
	}
	client := new(http.Client)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accepts", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	var data = ReleaseJson{}
	if response.StatusCode == http.StatusOK {
		json.Unmarshal(body, &data)
		return data.TagName, nil
	}

	return "", err
}

// convertVersion takes a string like "v1.2.3" and returns it
// as a number that can be accurately compared
func convertVersion(input string) int {
	cleanNumber := strings.TrimPrefix(input, "v")
	parts := strings.Split(cleanNumber, ".")

	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	patch, _ := strconv.Atoi(parts[2])

	return major*10000 + minor*100 + patch
}

// isOutdated returns true if the parsed local version string added
// to the binary is below the GitHub API response
func isOutdated(local string, published string) bool {
	latestVersion := convertVersion(published)
	installedVersion := convertVersion(local)

	return installedVersion < latestVersion
}

func CheckForVersion(ctx *AppContext) {
	tagName, err := requestPublishedVersion()
	if err != nil {
		fmt.Println(fmt.Errorf("err %v", err))
	}

	if isOutdated(ctx.Version, tagName) {
		title := fmt.Sprintf("A new version (%s) is available.", tagName)
		titleEl := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6")).Render(title)
		fmt.Println(box.Render(heredoc.Docf(`
      %s

      You have %s installed.
      Please download the latest version to continue access to our services.
    `, titleEl, ctx.Version)))
	}
}
