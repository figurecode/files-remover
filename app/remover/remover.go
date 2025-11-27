package remover

import (
	"io"
	"os"
	"text/template"

	"github.com/figurecode/files-remover/conf"
	"github.com/figurecode/files-remover/scanner"
)

const debugReportTempl = `{{.FilesCount}} files will be deleted in total
{{.Size}} will be freed up

Files to be deleted:
{{range .Files}}---------------------------------
PATH: {{.}}
{{end}}`

type Remover interface {
	Execute(files scanner.FoundFiles) error
}

type DebugRemover struct {
	outStream io.Writer
}

func (d DebugRemover) Execute(files scanner.FoundFiles) error {
	if len(files) == 0 {
		files = make(scanner.FoundFiles)
	}

	var reportParam struct {
		FilesCount int
		Files      []string
		Size       int64
	}
	var report = template.Must(template.New("Debug mode").Parse(debugReportTempl))

	reportParam.FilesCount = len(files)

	for path, size := range files {
		reportParam.Files = append(reportParam.Files, path)
		reportParam.Size += size
	}

	if err := report.Execute(d.outStream, reportParam); err != nil {
		return err
	}

	return nil
}

type ActionRemover struct{}

func (a ActionRemover) Execute(files scanner.FoundFiles) error {
	if len(files) == 0 {
		return nil
	}

	for path := range files {
		err := os.Remove(path)

		if !os.IsNotExist(err) && err != nil {
			return err
		}
	}

	return nil
}

func NewRemover(cfg conf.Config) Remover {
	if cfg.IsDemo {
		return DebugRemover{
			outStream: cfg.OutStream,
		}
	}

	return ActionRemover{}
}
