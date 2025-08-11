package original_file_names

import (
	"github.com/PocketBuilds/xpb"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
)

func init() {
	xpb.Register(&Plugin{})
}

type Plugin struct{}

func (p *Plugin) Name() string {
	return "original_file_names"
}

var version string

func (p *Plugin) Version() string {
	return version
}

func (p *Plugin) Description() string {
	return "Keep original file names."
}

func (p *Plugin) Init(app core.App) error {
	app.OnRecordUpdateRequest().Bind(
		&hook.Handler[*core.RecordRequestEvent]{
			Func:     p.keepOriginalFilenames,
			Priority: -999,
		},
	)
	app.OnRecordCreateRequest().Bind(
		&hook.Handler[*core.RecordRequestEvent]{
			Func:     p.keepOriginalFilenames,
			Priority: -999,
		},
	)
	return nil
}

// https://github.com/pocketbase/pocketbase/discussions/2787
func (p *Plugin) keepOriginalFilenames(e *core.RecordRequestEvent) error {
	for _, field := range e.Collection.Fields {
		if field.Type() != core.FieldTypeFile {
			continue
		}
		keys := []string{field.GetId(), field.GetName()}
		for _, key := range keys {
			files := e.Record.GetUploadedFiles(key)
			for _, file := range files {
				file.Name = file.OriginalName
			}
			e.Record.Set(key, files)
		}
	}
	return e.Next()
}
