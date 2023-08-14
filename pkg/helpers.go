package pkg

import "strings"

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
