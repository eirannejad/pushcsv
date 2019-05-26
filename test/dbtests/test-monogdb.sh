#!
echo "running pushcsv tests on mongodb..."

# mongodb root account
ROOTU="mdbroot"
ROOTP="mdbrootpass"

# mongodb user for this test
TESTU="pushcsvtestuser"
TESTP="testpass"

# database and collection for this test
TESTDB="pushcsvtestdb"
TESTCOL="testcollection"

printf $PROGRESSPRINT "pulling docker images and starting a container..."
# pull docker container and create the root user
docker run -d  --name pushcsvtestmongo \
               -p 127.0.0.1:27017:27017/tcp \
               -e MONGO_INITDB_ROOT_USERNAME=$ROOTU -e MONGO_INITDB_ROOT_PASSWORD=$ROOTP mongo \
                 >>$LOGFILE 2>&1 
# wait for container to start
sleep 2

printf $PROGRESSPRINT "configuring test db..."
# report mongodb version
echo "mongo --version" | docker exec -i pushcsvtestmongo bash   >>$LOGFILE 2>&1 
sleep 2

# switch to test db and add test user
echo "use $TESTDB
db.createUser({user: \"$TESTU\", pwd: \"$TESTP\", roles: [\"dbOwner\"]})
exit
" | docker exec -i pushcsvtestmongo mongo -u $ROOTU -p $ROOTP   >>$LOGFILE 2>&1 
sleep 2

printf $PROGRESSPRINT
# test pushing utf-8 csv input
TEST="utf-8 test"
pushcsv mongodb://$TESTU:$TESTP@localhost:27017/$TESTDB $TESTCOL $ASSETS/test-utf8.csv --headers --debug --trace  >>$LOGFILE 2>&1 
if [ $? -eq 0 ]; then
    echo -e $PASS $TEST
else
    echo -e $FAIL $TEST
fi

# test pushing utf-8 with BOM csv input
TEST="utf-8 BOM test"
pushcsv mongodb://$TESTU:$TESTP@localhost:27017/$TESTDB $TESTCOL $ASSETS/test-utf8bom.csv --headers --debug --trace >>$LOGFILE 2>&1 
if [ $? -eq 0 ]; then
    echo -e $PASS $TEST
else
    echo -e $FAIL $TEST
fi

# test pushing utf-16 csv input
TEST="utf-16 test"
pushcsv mongodb://$TESTU:$TESTP@localhost:27017/$TESTDB $TESTCOL $ASSETS/test-utf16.csv --headers --debug --trace  >>$LOGFILE 2>&1 
if [ $? -eq 0 ]; then
    echo -e $PASS $TEST
else
    echo -e $FAIL $TEST
fi

# test pushing utf-16 with BOM csv input
TEST="utf-16 BOM test"
pushcsv mongodb://$TESTU:$TESTP@localhost:27017/$TESTDB $TESTCOL $ASSETS/test-utf16be.csv --headers --debug --trace  >>$LOGFILE 2>&1 
if [ $? -eq 0 ]; then
    echo -e $PASS $TEST
else
    echo -e $FAIL $TEST
fi

printf $PROGRESSPRINT "cleaning up..."
docker stop pushcsvtestmongo   >>$LOGFILE 2>&1 
docker rm pushcsvtestmongo  >>$LOGFILE 2>&1 
docker rmi mongo >>$LOGFILE 2>&1 