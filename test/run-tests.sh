#! /bin/bash
# globals
export TESTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
export TESTENV=$TESTPATH/env
export DBTESTS=$TESTPATH/dbs
export TESTCASES=$TESTPATH/cases
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
$TESTENV/setup-env.sh      | tee -a $LOGFILE

# db tests
$DBTESTS/test-postgres.sh  | tee -a $LOGFILE
$DBTESTS/test-monogdb.sh   | tee -a $LOGFILE
$DBTESTS/test-mysql.sh     | tee -a $LOGFILE
$DBTESTS/test-sqlserver.sh | tee -a $LOGFILE
$DBTESTS/test-sqlite.sh    | tee -a $LOGFILE

# teardown
$TESTENV/teardown-env.sh   | tee -a $LOGFILE
