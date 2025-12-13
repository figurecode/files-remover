package remover

import (
	"fmt"
	"io"
	"os"
	"text/template"
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

func DebugRemover(files map[string]int64, out io.Writer) error {
	if len(files) == 0 {
		files = make(map[string]int64)
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
	reportParam.Files = make([]string, 0, len(files))

	for path, size := range files {
		reportParam.Files = append(reportParam.Files, path)
		reportParam.Size += size
	}

	if err := report.Execute(out, reportParam); err != nil {
		return err
	}

	return nil
}

func Execute(files map[string]int64) error {
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
