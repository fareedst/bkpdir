module bkpdir

go 1.21

require (
	bkpdir/pkg/fileops v0.0.0
	bkpdir/pkg/formatter v0.0.0
	github.com/BurntSushi/toml v1.5.0
	github.com/spf13/cobra v1.8.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/bmatcuk/doublestar/v4 v4.8.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace bkpdir/pkg/fileops => ./pkg/fileops

replace bkpdir/pkg/formatter => ./pkg/formatter
