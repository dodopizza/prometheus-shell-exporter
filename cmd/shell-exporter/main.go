package main

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
)

// sanitizePromLabelName -
func sanitizePromLabelName(str string) string {
	re := regexp.MustCompile(`[\.\-]`)
	result := re.ReplaceAllString(str, "_")
	re = regexp.MustCompile(`^\d`)
	result = re.ReplaceAllString(result, "_$0")
	return result
}

// getScriptsInDir -
func getScriptsInDir(dir string, pattern string) (fNames []string, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			fNames = append(fNames, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fNames, nil
}

func main() {

	scripts, err := getScriptsInDir("/workspaces/prometheus-shell-exporter/examples", "*.json")
	if err != nil {
		errors.New(err.Error())
	}

	pe := NewPromExporter()

	for _, script := range scripts {
		scriptBaseName := filepath.Base(script)
		metricName := scriptBaseName[0 : len(scriptBaseName)-len(filepath.Ext(scriptBaseName))]
		p := PromMetrics{}
		p.ReadFromFile(script)
		pe.NewGaugeVecFromPromMetrics(sanitizePromLabelName(metricName), p)
	}

	err = pe.Serve()
	if err != nil {
		errors.New(err.Error())
	}

	os.Exit(0)
}
