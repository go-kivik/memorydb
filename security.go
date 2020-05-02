package memorydb

import (
	"context"
	"net/http"

	"github.com/go-kivik/kivik/v4"
	"github.com/go-kivik/kivik/v4/driver"
)

func cloneSecurity(in *driver.Security) *driver.Security {
	return &driver.Security{
		Admins: driver.Members{
			Names: in.Admins.Names,
			Roles: in.Admins.Roles,
		},
		Members: driver.Members{
			Names: in.Members.Names,
			Roles: in.Members.Roles,
		},
	}
}

func (d *db) Security(_ context.Context) (*driver.Security, error) {
	d.db.mu.RLock()
	defer d.db.mu.RUnlock()
	if d.db.deleted {
		return nil, &kivik.Error{HTTPStatus: http.StatusNotFound, Message: "missing"}
	}
	return cloneSecurity(d.db.security), nil
}

func (d *db) SetSecurity(_ context.Context, sec *driver.Security) error {
	d.db.mu.Lock()
	defer d.db.mu.Unlock()
	if d.db.deleted {
		return &kivik.Error{HTTPStatus: http.StatusNotFound, Message: "missing"}
	}
	d.db.security = cloneSecurity(sec)
	return nil
}
