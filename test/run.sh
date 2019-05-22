#!
TESTPATH="`dirname \"$0\"`"

echo "# setting up test environment..."
$TESTPATH/setup-env.sh

echo "# testing postgresql..."
$TESTPATH/test-postgres.sh

echo "# testing mongodb..."
$TESTPATH/test-monogdb.sh

echo "# testing mysql..."
$TESTPATH/test-mysql.sh

echo "# testing sql server..."
$TESTPATH/test-sqlserver.sh

echo "# testing sqlite3..."
$TESTPATH/test-sqlite.sh

echo "# tearing down test environment..."
$TESTPATH/teardown-env.sh
