#!
TESTPATH="`dirname \"$0\"`"
TABLENAME="sometable"

# pushcsv sqlite3:data.db $TABLENAME $TESTPATH/test-utf8.csv --dry-run --debug --trace 2>&1 >/dev/null | grep -E "encoding|INSERT"
echo "testing sqlite3:utf8..."
pushcsv sqlite3:data.db $TABLENAME $TESTPATH/test-utf8.csv --dry-run --debug --trace 2>&1 >/dev/null

echo "testing sqlite3:utf8bom..."
pushcsv sqlite3:data.db $TABLENAME $TESTPATH/test-utf8bom.csv --dry-run --debug --trace 2>&1 >/dev/null

echo "testing sqlite3:utf16..."
pushcsv sqlite3:data.db $TABLENAME $TESTPATH/test-utf16.csv --dry-run --debug --trace 2>&1 >/dev/null

echo "testing sqlite3:utf16be..."
pushcsv sqlite3:data.db $TABLENAME $TESTPATH/test-utf16be.csv --dry-run --debug --trace 2>&1 >/dev/null

echo "cleaning up sqlite3 datastore..."
rm data.db