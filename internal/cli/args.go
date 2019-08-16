package cli

import (
	"strings"

	"github.com/docopt/docopt-go"
)

type Options struct {
	Opts       *docopt.Opts
	ConnString string
	Table      string
	DataFile   string
	HasHeaders bool
	Purge      bool
	AttrMaps   map[string]string
	Debug      bool
	Trace      bool
	DryRun     bool
	Compact    bool
}

func NewOptions(argv []string) *Options {

	parser := &docopt.Parser{
		HelpHandler: printHelpAndExit,
	}

	opts, _ := parser.ParseArgs(help, argv, version)

	connString, _ := opts.String("<db_uri>")
	table, _ := opts.String("<table>")
	dataFile, _ := opts.String("<file>")
	hasHeaders, _ := opts.Bool("--headers")
	purge, _ := opts.Bool("--purge")

	// --map is a repeated argument and value is of type []string but
	// passed as generic interface{} so needs type assertion
	attrmapArg, _ := opts["--map"].([]string)
	attrMaps := make(map[string]string)
	for _, attrMapStr := range attrmapArg {
		parts := strings.Split(attrMapStr, ":")
		attrMaps[parts[0]] = parts[1]
	}

	// first line must be treated as header if mapping is specified
	if len(attrMaps) > 0 {
		hasHeaders = true
	}

	debug, _ := opts.Bool("--debug")
	trace, _ := opts.Bool("--trace")
	dryRun, _ := opts.Bool("--dry-run")
	compact, _ := opts.Bool("--compact")

	return &Options{
		Opts:       &opts,
		ConnString: connString,
		Table:      table,
		DataFile:   dataFile,
		HasHeaders: hasHeaders,
		Purge:      purge,
		AttrMaps:   attrMaps,
		Debug:      debug,
		Trace:      trace,
		DryRun:     dryRun,
		Compact:    compact,
	}
}
