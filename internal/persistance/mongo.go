package persistance

import (
	"fmt"

	"github.com/eirannejad/pushcsv/internal/datafile"
	_ "github.com/lib/pq"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDBWriter struct {
	DatabaseWriter
}

func (w MongoDBWriter) Write(tableData *datafile.TableData) (*Result, error) {
	// parse and grab database name from uri
	w.Logger.Debug("grabbing db name from connection string")
	dialinfo, err := mgo.ParseURL(w.Config.ConnString)
	if err != nil {
		return nil, err
	}

	w.Logger.Debug("opening mongodb session")
	session, cErr := mgo.DialWithInfo(dialinfo)
	if cErr != nil {
		return nil, cErr
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	w.Logger.Debug("setting session properties")
	session.SetMode(mgo.Monotonic, true)
	w.Logger.Trace(session)

	w.Logger.Debug("getting target collection")
	c := session.DB(dialinfo.Database).C(tableData.Name)
	w.Logger.Trace(c)

	if len(tableData.Records) > 0 {
		if tableData.HasHeaders() {
			// build sql data info
			var count int
			w.Logger.Debug("opening bulk operation")
			bulkop := c.Bulk()

			// purge the collection if requested
			if w.Purge {
				w.Logger.Debug("purging all existing documents")
				bulkop.RemoveAll(bson.M{})
			}

			w.Logger.Debug("building documents")
			for _, record := range tableData.Records {
				map_obj := make(map[string]string)
				for fidx, field := range record {
					map_obj[tableData.Headers[fidx]] = field
				}
				w.Logger.Trace(map_obj)
				bulkop.Insert(map_obj)
				count++
			}
			if !w.DryRun {
				w.Logger.Debug("running bulk operation")
				_, txnErr := bulkop.Run()
				if txnErr != nil {
					return nil, txnErr
				}
				w.Logger.Debug("preparing report")
				return &Result{
					Message: fmt.Sprintf("successfully updated %d documents", count),
				}, nil
			} else {
				w.Logger.Debug("dry run complete")
				return &Result{
					ResultCode: 2,
					Message:    "processed documents but no changed were made to db",
				}, nil
			}
		} else {
			return &Result{
				ResultCode: 3,
				Message:    "headers are required",
			}, nil
		}
	}
	w.Logger.Debug("nothing to write")
	return &Result{
		ResultCode: 1,
		Message:    "no data to write",
	}, nil
}
