package shellexporter

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
)

func getMetricsScripts(scriptsDir string) (scripts []string, err error) {
	scriptDirAbs, err := filepath.Abs(scriptsDir)
	if err != nil {
		return
	}

	for _, ext := range []string{"*.ps1", "*.sh"} {
		s, err := walkMatch(scriptDirAbs, ext)
		if err != nil {
			return []string{}, err
		}
		scripts = append(scripts, s...)

	}

	if len(scripts) <= 0 {
		err = errors.New("no scripts to serve")
		return
	}

	return
}

func walkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func getFileName(fname string) string {
	fname = filepath.Base(fname)
	return fname[0 : len(fname)-len(filepath.Ext(fname))]
}

func sanitizePromLabelName(str string) string {
	re := regexp.MustCompile(`[\.\-]`)
	result := re.ReplaceAllString(str, "_")
	re = regexp.MustCompile(`^\d`)
	result = re.ReplaceAllString(result, "_$0")
	return result
}
