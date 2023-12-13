package pkg

import "strings"

// ArgumentsSliceToString takes a cobra supplied slice of strings and joins them
// based on the total length. For example:
//
//	args := []string{"api_key", "host", "namespace"}
//	argsSentence := pkg.ArgumentsSliceToString(args, "and")
//
// prints "api_key, host and namespace"
func ArgumentsSliceToString(args []string, separator string) string {
	argsLength := len(args)
	switch argsLength {
	case 0:
		return ""
	case 1:
		return args[0]
	case 2:
		return args[0] + " " + separator + " " + args[1]
	default:
		return strings.Join(args[:argsLength-1], ", ") + " " + separator + " " + args[argsLength-1]
	}
}
