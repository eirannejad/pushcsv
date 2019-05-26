# all test cases for pushcsv
# TESTCONFIG is setup with pushcsv arguments (db rui and table) by each
# database test script
# -h argument requests to use --headers argument on test cases e.g. mongodb
# tests need headers for object property names

# report result of last command
# $1 is test name
function report() {
    if [ $? -eq 0 ]; then
        echo -e $PASS $1
    else
        echo -e $FAIL $1
    fi
}

export -f report

# run tests
$TESTCASES/test-encoding.sh $*
$TESTCASES/test-pushall.sh  $*
$TESTCASES/test-pushmap.sh  $*
$TESTCASES/test-purge.sh    $*