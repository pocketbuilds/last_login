package last_login

import (
	"errors"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/pocketbuilds/xpb"
)

type Plugin struct {
	// Datetime fields on auth collections that will
	//   automatically set on record auth request.
	//   - format "<collection_name>.<field_name>"
	Fields []string `json:"fields"`
}

func init() {
	xpb.Register(&Plugin{})
}

// Name implements xpb.Plugin
func (p *Plugin) Name() string {
	return "last_login"
}

// This variable will automatically be set at build time by xpb
var version string

// Version implements xpb.Plugin
func (p *Plugin) Version() string {
	return version
}

// Description implements xpb.Plugin
func (p *Plugin) Description() string {
	return "Allow for easily configured last_login fields."
}

// Init implements xpb.Plugin
func (p *Plugin) Init(app core.App) error {
	for _, raw := range p.Fields {
		collectionName, fieldName, _ := strings.Cut(raw, ".")
		app.OnRecordAuthRequest(collectionName).
			BindFunc(p.setLastLoginField(fieldName))
	}
	return nil
}

// copied from core.collectionNameRegex
var collectionNameRegex = regexp.MustCompile(`^\w+$`)

// Validate implements validation.Validatable.
func (p *Plugin) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Fields,
			validation.Each(
				validation.By(func(value any) error {
					str, err := validation.EnsureString(value)
					if err != nil {
						return err
					}
					if strings.Count(str, ".") != 1 {
						return errors.New("must contain a single period separator")
					}
					collectionName, fieldName, _ := strings.Cut(str, ".")
					if err := validation.Validate(collectionName,
						validation.Match(collectionNameRegex),
					); err != nil {
						return err
					}
					if err := validation.Validate(fieldName,
						validation.By(core.DefaultFieldNameValidationRule),
					); err != nil {
						return err
					}
					return nil
				}),
			),
		),
	)
}

func (p *Plugin) setLastLoginField(fieldName string) func(e *core.RecordAuthRequestEvent) error {
	return func(e *core.RecordAuthRequestEvent) error {
		e.Record.Set(fieldName, types.NowDateTime())
		if err := e.App.Save(e.Record); err != nil {
			return err
		}
		return e.Next()
	}
}
