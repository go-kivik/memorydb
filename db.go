package memorydb

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-kivik/kivik/v4"
	"github.com/go-kivik/kivik/v4/driver"
)

var notYetImplemented = &kivik.Error{Status: http.StatusNotImplemented, Message: "kivik: not yet implemented in memory driver"}

// database is an in-memory database representation.
type db struct {
	*client
	dbName string
	db     *database
}

func (d *db) Query(ctx context.Context, ddoc, view string, opts map[string]interface{}) (driver.Rows, error) {
	// FIXME: Unimplemented
	return nil, notYetImplemented
}

func (d *db) Get(ctx context.Context, docID string, opts map[string]interface{}) (*driver.Document, error) {
	if exists, _ := d.client.DBExists(ctx, d.dbName, nil); !exists {
		return nil, &kivik.Error{Status: http.StatusPreconditionFailed, Message: "database does not exist"}
	}
	if !d.db.docExists(docID) {
		return nil, &kivik.Error{Status: http.StatusNotFound, Message: "missing"}
	}
	if rev, ok := opts["rev"].(string); ok {
		if doc, found := d.db.getRevision(docID, rev); found {
			return &driver.Document{
				ContentLength: int64(len(doc.data)),
				Rev:           rev,
				Body:          ioutil.NopCloser(bytes.NewReader(doc.data)),
			}, nil
		}
		return nil, &kivik.Error{Status: http.StatusNotFound, Message: "missing"}
	}
	last, _ := d.db.latestRevision(docID)
	if last.Deleted {
		return nil, &kivik.Error{Status: http.StatusNotFound, Message: "missing"}
	}
	return &driver.Document{
		ContentLength: int64(len(last.data)),
		Rev:           fmt.Sprintf("%d-%s", last.ID, last.Rev),
		Body:          ioutil.NopCloser(bytes.NewReader(last.data)),
	}, nil
}

func (d *db) CreateDoc(ctx context.Context, doc interface{}, _ map[string]interface{}) (docID, rev string, err error) {
	if exists, _ := d.client.DBExists(ctx, d.dbName, nil); !exists {
		return "", "", &kivik.Error{Status: http.StatusPreconditionFailed, Message: "database does not exist"}
	}
	couchDoc, err := toCouchDoc(doc)
	if err != nil {
		return "", "", err
	}
	if id, ok := couchDoc["_id"].(string); ok {
		docID = id
	} else {
		docID = randStr()
	}
	rev, err = d.Put(ctx, docID, doc, nil)
	return docID, rev, err
}

func (d *db) Put(ctx context.Context, docID string, doc interface{}, _ map[string]interface{}) (rev string, err error) {
	if exists, _ := d.client.DBExists(ctx, d.dbName, nil); !exists {
		return "", &kivik.Error{Status: http.StatusPreconditionFailed, Message: "database does not exist"}
	}
	isLocal := strings.HasPrefix(docID, "_local/")
	if !isLocal && docID[0] == '_' && !strings.HasPrefix(docID, "_design/") {
		return "", &kivik.Error{Status: http.StatusBadRequest, Message: "Only reserved document ids may start with underscore."}
	}
	couchDoc, err := toCouchDoc(doc)
	if err != nil {
		return "", err
	}
	couchDoc["_id"] = docID
	// TODO: Add support for storing attachments.
	delete(couchDoc, "_attachments")

	if last, ok := d.db.latestRevision(docID); ok {
		if !last.Deleted && !isLocal && couchDoc.Rev() != fmt.Sprintf("%d-%s", last.ID, last.Rev) {
			return "", &kivik.Error{Status: http.StatusConflict, Message: "document update conflict"}
		}
		return d.db.addRevision(couchDoc), nil
	}

	if couchDoc.Rev() != "" {
		// Rev should not be set for a new document
		return "", &kivik.Error{Status: http.StatusConflict, Message: "document update conflict"}
	}
	return d.db.addRevision(couchDoc), nil
}

var revRE = regexp.MustCompile("^[0-9]+-[a-f0-9]{32}$")

func validRev(rev string) bool {
	return revRE.MatchString(rev)
}

func (d *db) Delete(ctx context.Context, docID, rev string, _ map[string]interface{}) (newRev string, err error) {
	if exists, _ := d.client.DBExists(ctx, d.dbName, nil); !exists {
		return "", &kivik.Error{Status: http.StatusPreconditionFailed, Message: "database does not exist"}
	}
	if !strings.HasPrefix(docID, "_local/") && !validRev(rev) {
		return "", &kivik.Error{Status: http.StatusBadRequest, Message: "invalid rev format"}
	}
	if !d.db.docExists(docID) {
		return "", &kivik.Error{Status: http.StatusNotFound, Message: "missing"}
	}
	return d.Put(ctx, docID, map[string]interface{}{
		"_id":      docID,
		"_rev":     rev,
		"_deleted": true,
	}, nil)
}

func (d *db) Stats(_ context.Context) (*driver.DBStats, error) {
	return &driver.DBStats{
		Name: d.dbName,
		// DocCount:     0,
		// DeletedCount: 0,
		// UpdateSeq:    "",
		// DiskSize:     0,
		// ActiveSize:   0,
		// ExternalSize: 0,
	}, nil
}

func (c *client) Compact(_ context.Context) error {
	// FIXME: Unimplemented
	return notYetImplemented
}

func (d *db) CompactView(_ context.Context, _ string) error {
	// FIXME: Unimplemented
	return notYetImplemented
}

func (d *db) ViewCleanup(_ context.Context) error {
	// FIXME: Unimplemented
	return notYetImplemented
}

func (d *db) Changes(ctx context.Context, opts map[string]interface{}) (driver.Changes, error) {
	// FIXME: Unimplemented
	return nil, notYetImplemented
}

func (d *db) PutAttachment(_ context.Context, _, _ string, _ *driver.Attachment, _ map[string]interface{}) (string, error) {
	// FIXME: Unimplemented
	return "", notYetImplemented
}

func (d *db) GetAttachment(ctx context.Context, docID, filename string, opts map[string]interface{}) (*driver.Attachment, error) {
	// FIXME: Unimplemented
	return nil, notYetImplemented
}

func (d *db) DeleteAttachment(ctx context.Context, docID, rev, filename string, opts map[string]interface{}) (newRev string, err error) {
	// FIXME: Unimplemented
	return "", notYetImplemented
}
