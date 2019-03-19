package persistance

import (
	"github.com/pkg/errors"
	// "database/sql"
	// "fmt"
	"github.com/eirannejad/pushcsv/internal/csv"
	_ "github.com/lib/pq"
	// "strings"
)

type MongoDBWriter struct {
	DatabaseWriter
}

// TODO Update
func (w MongoDBWriter) Write(csvData *csv.CsvData) (*Result, error) {
	return nil, errors.New("mongodb interface not yet implemented.")
}

// 	// parse and grab database name from uri
// 	dialinfo, err := mgo.ParseURL(connstr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	dbname := dialinfo.Database

// 	// connect to db engine
// 	session, err := mgo.Dial(connstr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer session.Close()
// 	// Optional. Switch the session to a monotonic behavior.
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB(dbname).C(table)

// 	// purge the collection if requested
// 	if purge {
// 		c.RemoveAll(bson.M{})
// 	}

// 	if len(docs) > 0 {
// 		if len(attrs) == 0 {
// 			if len(headers) == 0 {
// 				log.Fatal("`--map` must be specified when pushing a csv with no headers.")
// 				return 0, nil
// 			} else {
// 				attrs = headers
// 			}
// 		}

// 		// build sql data info
// 		bulkop := c.Bulk()
// 		for _, record := range docs {
// 			map_obj := make(map[string]string)
// 			for fidx, field := range record {
// 				map_obj[attrs[fidx]] = field
// 			}
// 			bulkop.Insert(map_obj)
// 		}
// 		res, err := bulkop.Run()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		return res.Modified, nil
// 	}
// 	return 0, nil
// }
