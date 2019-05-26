#!
echo "running pushcsv tests on sqlite3..."

TABLENAME="sqlitetesttable"
TESTDB=$TESTPATH/sqlitetestdb.db

TEST="utf-8 test"
pushcsv sqlite3:$TESTDB $TABLENAME $ASSETS/test-utf8.csv --debug --trace >>$LOGFILE 2>&1 
if [ $? -eq 0 ]; then
    echo -e $PASS $TEST
else
    echo -e $FAIL $TEST
fi

TEST="utf-8 BOM test"
pushcsv sqlite3:$TESTDB $TABLENAME $ASSETS/test-utf8bom.csv --debug --trace >>$LOGFILE 2>&1 
if [ $? -eq 0 ]; then
    echo -e $PASS $TEST
else
    echo -e $FAIL $TEST
fi

TEST="utf-16 test"
pushcsv sqlite3:$TESTDB $TABLENAME $ASSETS/test-utf16.csv --debug --trace >>$LOGFILE 2>&1 
if [ $? -eq 0 ]; then
    echo -e $PASS $TEST
else
    echo -e $FAIL $TEST
fi

TEST="utf-16 BOM test"
pushcsv sqlite3:$TESTDB $TABLENAME $ASSETS/test-utf16be.csv --debug --trace >>$LOGFILE 2>&1 
if [ $? -eq 0 ]; then
    echo -e $PASS $TEST
else
    echo -e $FAIL $TEST
fi

printf $PROGRESSPRINT "cleaning up..."
rm $TESTDB
