package last_login

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/pocketbuilds/xpb"
)

func init() {
	xpb.Register(&Plugin{
		FieldName: "last_login",
	})
}

type Plugin struct {
	// The name given to the last_login field.
	// default: "last_login"
	FieldName string `json:"field_name"`
}

func (p *Plugin) Name() string {
	return "last_login"
}

var version string

func (p *Plugin) Version() string {
	return version
}

func (p *Plugin) Description() string {
	return "Adds a last_login datetime field to all auth collections, and updates it on each auth request."
}

func (p *Plugin) Init(app core.App) error {
	app.OnServe().BindFunc(p.migrateLastLoginField)
	app.OnRecordAuthRequest().BindFunc(p.setLastLogin)
	return nil
}

func (p *Plugin) migrateLastLoginField(e *core.ServeEvent) error {
	authCollections, err := e.App.FindAllCollections(core.CollectionTypeAuth)
	if err != nil {
		return err
	}
	err = e.App.RunInTransaction(func(txApp core.App) error {
		for _, collection := range authCollections {
			// Skip if field already exists and is of type schema.Date
			if field := collection.Fields.GetByName(p.FieldName); field != nil && field.Type() == core.FieldTypeDate {
				continue
			}
			collection.Fields.Add(&core.DateField{
				Name: p.FieldName,
			})

			if err := txApp.Save(collection); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return e.Next()
}

func (p *Plugin) setLastLogin(e *core.RecordAuthRequestEvent) error {
	e.Record.Set(p.FieldName, types.NowDateTime())
	if err := e.App.Save(e.Record); err != nil {
		return err
	}
	return e.Next()
}
