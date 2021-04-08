package shellexporter

import (
	"testing"
)

func Test_sanitizePromLabelName_ShouldReturnValidPrometheusLabelName(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"NoDots", args{"Label.test"}, "Label_test"},
		{"NoDashes", args{"Label-test-two"}, "Label_test_two"},
		{"Number as first", args{"7zip"}, "_7zip"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitizePromLabelName(tt.args.str); got != tt.want {
				t.Errorf("sanitizePromLabelName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFileName_ShouldReturnFileNameWithoutExtension(t *testing.T) {
	type args struct {
		fname string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "File with extension", args: args{"Hello.exe"}, want: "Hello"},
		{name: "File without extension", args: args{"Hello"}, want: "Hello"},
		{name: "File path with extension", args: args{"/tmp/Hello"}, want: "Hello"},
		{name: "File with double extension", args: args{"Hello.tar.gz"}, want: "Hello.tar"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFileName(tt.args.fname); got != tt.want {
				t.Errorf("getFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
