package remover

import (
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/figurecode/files-remover/conf"
	"github.com/figurecode/files-remover/scanner"
)

const debugReportTempl = `{{.FilesCount}} files will be deleted in total
{{humanSize .Size}} of disk space will be freed

Files to be deleted:
{{range .Files}}---------------------------------
PATH: {{.}}
{{end}}
END
`

func humanSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

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
	var report = template.Must(
		template.New("Debug mode").
			Funcs(template.FuncMap{"humanSize": humanSize}).
			Parse(debugReportTempl))

	reportParam.FilesCount = len(files)
	reportParam.Files = make([]string, len(files))

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
