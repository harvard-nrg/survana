package dashboard

import (
	"neuroinformatics.harvard.edu/survana"
)

const (
	NAME = "dashboard"
)

//The Admin component
type Dashboard struct {
	*survana.Module
	mux *survana.RESTMux
}

// creates a new Admin module
func NewModule(path string, db survana.Database) *Dashboard {

	mux := survana.NewRESTMux()

	m := &Dashboard{
		Module: &survana.Module{
			Name:   NAME,
			Path:   path,
			Db:     db,
			Router: mux,
			Log:    db.NewLogger("logs", NAME),
		},
		mux: mux,
	}

	m.ParseTemplates()

	m.RegisterHandlers()

	return m
}