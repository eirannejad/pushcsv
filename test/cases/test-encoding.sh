echo "testing support for various csv encoding..."

TESTARGS="--debug --trace >>$LOGFILE 2>&1"
if [[ $3 == "-h" ]]; then
    TESTARGS="--headers "$TESTARGS
fi

# test pushing utf-8 csv input
eval pushcsv $1 $2 $ASSETS/test-utf8.csv $TESTARGS
report "utf-8 test"

# test pushing utf-8 with BOM csv input
eval pushcsv $1 $2 $ASSETS/test-utf8bom.csv $TESTARGS
report "utf-8 BOM test"

# test pushing utf-16 csv input
eval pushcsv $1 $2 $ASSETS/test-utf16.csv $TESTARGS
report "utf-16 test"

# test pushing utf-16 with BOM csv input
eval pushcsv $1 $2 $ASSETS/test-utf16be.csv $TESTARGS
report "utf-16 BOM test"