package clioutput

import "strings"

var textureTokens = map[string]struct{}{
	"--plain":    {},
	"-p":         {},
	"--no-color": {},
	"--stdout":   {},
}

// PeelTextureArgs removes global operator-output flags from argv before Cobra or auto-detect parsing.
func PeelTextureArgs(argv []string) (remaining []string, flags PeeledFlags) {
	for i := 0; i < len(argv); i++ {
		arg := argv[i]
		if _, ok := textureTokens[arg]; ok {
			switch arg {
			case "--plain", "-p":
				flags.Plain = true
			case "--no-color":
				flags.NoColor = true
			case "--stdout":
				flags.Stdout = true
			}
			continue
		}
		if strings.HasPrefix(arg, "--plain=") {
			flags.Plain = true
			continue
		}
		remaining = append(remaining, arg)
	}
	return remaining, flags
}
