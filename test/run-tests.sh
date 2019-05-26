#!
# globals
export TESTPATH="`dirname \"$0\"`"
export TESTENV=$TESTPATH/testenv
export DBTESTS=$TESTPATH/dbtests
export ASSETS=$TESTPATH/assets
export LOGPATH=$TESTPATH/logs

# green PASS
export PASS="[ \033[32mPASS\033[39m ]"
# red FAIL
export FAIL="[ \033[31mFAIL\033[39m ]"
export PROGRESSPRINT="%-120s\r"

# get HEAD commit hash for this test
export HEADHASH=`git describe --always`
export LOGFILE=$LOGPATH/$HEADHASH"_test.log"

# setup
$TESTENV/setup-env.sh

# db tests
$DBTESTS/test-postgres.sh
# $DBTESTS/test-monogdb.sh
$DBTESTS/test-mysql.sh
$DBTESTS/test-sqlserver.sh
$DBTESTS/test-sqlite.sh

# teardown
$TESTENV/teardown-env.sh
