package main

import "testing"

func Test_normalize(t *testing.T) {

	tests := []struct {
		name  string
		base  string
		file  string
		strip bool
		want  string
	}{
		{"stripBase", "/var/www/html/", "/var/www/html/test/a.gif", true, "test/a.gif"},
		{"don't stripBase", "/var/www/html", "/var/www/html/test/a.gif", false, "html/test/a.gif"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalize(tt.base, tt.file, tt.strip); got != tt.want {
				t.Errorf("stripBase() = %v, want %v", got, tt.want)
			}
		})
	}
}
