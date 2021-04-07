package main

import (
	"testing"
)

var _ = func() bool { // https://github.com/golang/go/issues/31859
	testing.Init()
	return true
}()

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
