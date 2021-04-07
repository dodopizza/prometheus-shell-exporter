package shellexporter

import (
	"os"
	"path/filepath"
	"regexp"
)

func WalkMatch(root, pattern string) ([]string, error) {
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

func GetFileName(fname string) string {
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
