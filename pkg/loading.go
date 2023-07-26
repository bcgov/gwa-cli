package pkg

import (
	"time"

	"github.com/briandowns/spinner"
)

func NewSpinner() *spinner.Spinner {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	return s
}
