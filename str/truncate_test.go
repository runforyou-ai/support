package str

import "testing"

func TestLimit(t *testing.T) {
	tests := []struct {
		s    string
		n    int
		end  string
		want string
	}{
		{"The quick brown fox", 9, "...", "The quick..."},
		{"The quick brown fox", 10, "...", "The quick..."},
		{"The quick brown fox", 19, "...", "The quick brown fox"},
		{"The quick brown fox", 100, "...", "The quick brown fox"},
		{"hello", 0, "...", "..."},
		{"hello", -1, "...", "..."},
		{"", 0, "...", ""},
		{"", -1, "...", ""},
		{"你好世界", 2, "…", "你好…"},
		{"👨\u200d👩x", 1, "…", "👨…"},
	}
	for _, tt := range tests {
		if got := Limit(tt.s, tt.n, tt.end); got != tt.want {
			t.Errorf("Limit(%q, %d, %q) = %q, want %q", tt.s, tt.n, tt.end, got, tt.want)
		}
	}
}

func TestWords(t *testing.T) {
	tests := []struct {
		s    string
		n    int
		end  string
		want string
	}{
		{"Perfectly balanced, as all things should be.", 3, " >>>", "Perfectly balanced, as >>>"},
		{"one two three", 3, "...", "one two three"},
		{"one two three  ", 3, "...", "one two three  "},
		{"one   two three", 1, "...", "one..."},
		{"  one two", 1, "...", "  one..."},
		{"one two", 0, "...", "..."},
		{"one two", -2, "...", "..."},
		{"", 2, "...", ""},
		{"你好 世界 再见", 2, "…", "你好 世界…"},
	}
	for _, tt := range tests {
		if got := Words(tt.s, tt.n, tt.end); got != tt.want {
			t.Errorf("Words(%q, %d, %q) = %q, want %q", tt.s, tt.n, tt.end, got, tt.want)
		}
	}
}

func TestExcerpt(t *testing.T) {
	tests := []struct {
		text, phrase string
		radius       int
		omission     string
		want         string
		wantOK       bool
	}{
		{"This is my name", "my", 3, "...", "...is my na...", true},
		{"This is my name", "name", 3, "...", "...my name", true},
		{"This is my name", "This", 3, "...", "This is...", true},
		{"This is my name", "MY", 100, "...", "This is my name", true},
		{"This is my name", "this", 0, "(...)", "This(...)", true},
		{"This is my name", "", 4, "...", "This...", true},
		{"", "", 4, "...", "", true},
		{"This is my name", "is", -1, "...", "...is...", true},
		{"This is my name", "nope", 3, "...", "", false},
		{"  padded  text  ", "padded", 100, "...", "padded  text", true},
		{"这是我的名字", "我的", 1, "…", "…是我的名…", true},
		{"", "", 2, "...", "", true},
	}
	for _, tt := range tests {
		got, ok := Excerpt(tt.text, tt.phrase, tt.radius, tt.omission)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("Excerpt(%q, %q, %d, %q) = %q, %v, want %q, %v",
				tt.text, tt.phrase, tt.radius, tt.omission, got, ok, tt.want, tt.wantOK)
		}
	}
}
