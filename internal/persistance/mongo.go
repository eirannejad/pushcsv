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

func (w MongoDBWriter) Purge(tableName string) (*Result, error) {
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
	c := session.DB(dialinfo.Database).C(tableName)
	w.Logger.Trace(c)

	// purge the collection if requested
	var purged int
	w.Logger.Debug("purging all existing documents")
	if !w.DryRun {
		cinfo, runErr := c.RemoveAll(bson.M{})
		if runErr != nil {
			return nil, runErr
		}
		purged = cinfo.Removed
	}

	w.Logger.Debug("preparing report")
	return &Result{
		Message: fmt.Sprintf("successfully purged %d records", purged),
	}, nil
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
	db := session.DB(dialinfo.Database)
	c := db.C(tableData.Name)
	w.Logger.Trace(c)

	w.Logger.Debug("opening bulk operation")
	bulkop := c.Bulk()

	// purge the collection if requested
	if w.PurgeBeforeWrite {
		w.Logger.Debug("purging all existing documents")
		if !w.DryRun {
			bulkop.RemoveAll(bson.M{})
		}
	}

	// build sql data info
	w.Logger.Debug("building documents")
	for _, record := range tableData.Records {
		map_obj := make(map[string]string)
		for fidx, field := range record {
			map_obj[tableData.Headers[fidx]] = field
		}
		w.Logger.Trace(map_obj)
		bulkop.Insert(map_obj)
	}

	w.Logger.Debug("running bulk operation")
	if !w.DryRun {
		_, txnErr := bulkop.Run()
		if txnErr != nil {
			return nil, txnErr
		}

		// compact collection if requested
		if w.CompactAfterWrite {
			w.Logger.Debug("compacting collection")
			if !w.DryRun {
				compactErr := db.Run(
					bson.D{{Name: "compact", Value: tableData.Name}},
					nil)
				if compactErr != nil {
					return nil, compactErr
				}
				w.Logger.Debug("not implemented. need compacting support in mongodb driver")
			}
		}

		w.Logger.Debug("preparing report")
		return &Result{
			Message: fmt.Sprintf(
				"successfully updated %d documents",
				len(tableData.Records)),
		}, nil
	} else {
		return &Result{
			ResultCode: 2,
			Message:    "processed documents but no changed were made to db",
		}, nil
	}
}
