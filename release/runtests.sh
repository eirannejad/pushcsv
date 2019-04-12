#!
pushcsv sqlite3:data.db sometable ./release/test-utf8.csv --dry-run --debug --trace 2>&1 >/dev/null | grep -E "encoding|INSERT"
pushcsv sqlite3:data.db sometable ./release/test-utf8bom.csv --dry-run --debug --trace 2>&1 >/dev/null | grep -E "encoding|INSERT"
pushcsv sqlite3:data.db sometable ./release/test-utf16.csv --dry-run --debug --trace 2>&1 >/dev/null | grep -E "encoding|INSERT"
pushcsv sqlite3:data.db sometable ./release/test-utf16be.csv --dry-run --debug --trace 2>&1 >/dev/null | grep -E "encoding|INSERT"
rm data.db