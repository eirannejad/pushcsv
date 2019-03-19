package cli

import (
	"github.com/docopt/docopt-go"
	"strings"
)

type Options struct {
	ConnString string
	Table      string
	CsvFile    string
	HasHeaders bool
	Purge      bool
	AttrMap    []string
	Debug      bool
}

func NewOptions(argv []string) *Options {

	parser := &docopt.Parser{
		HelpHandler: PrintHelpAndExit,
	}

	opts, _ := parser.ParseArgs(help, argv, version)

	connString, _ := opts.String("<db_uri>")
	table, _ := opts.String("<table>")
	csvFile, _ := opts.String("<csv_file>")
	hasHeaders, _ := opts.Bool("--headers")
	purge, _ := opts.Bool("--purge")
	attrmapArg, _ := opts.String("--map")
	debug, _ := opts.Bool("--debug")

	var attrmap []string
	if attrmapArg != "" {
		attrmap = strings.Split(attrmapArg, ";")
	} else {
		attrmap = make([]string, 0)
	}
	return &Options{
		ConnString: connString,
		Table:      table,
		CsvFile:    csvFile,
		HasHeaders: hasHeaders,
		Purge:      purge,
		AttrMap:    attrmap,
		Debug:      debug,
	}
}
