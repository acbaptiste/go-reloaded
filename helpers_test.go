package main

import (
	"math"
	"reflect"
	"testing"
)

func TestApplyMarkers(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"so exciting (up, 2)", []string{"SO", "EXCITING"}},
		{"word(cap)!", []string{"Word", "!"}},
		{"ff ( hex ) end", []string{"255", "end"}},
		{"hello (rev)", []string{"hello", "(rev)"}},
		{"(up) hello", []string{"hello"}},
		{"   ", nil},
	}
	for _, c := range cases {
		got := applyMarkers(c.in)
		if len(got) == 0 && len(c.want) == 0 {
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("applyMarkers(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestMarkerCount(t *testing.T) {
	cases := map[string]int{
		"":                        1,
		"6":                       6,
		"0":                       0,
		"99999999999999999999999": math.MaxInt,
	}
	for in, want := range cases {
		if got := markerCount(in); got != want {
			t.Errorf("markerCount(%q) = %d, want %d", in, got, want)
		}
	}
}

func TestConvertToDecimal(t *testing.T) {
	cases := []struct {
		word string
		base int
		want string
	}{
		{"ff", 16, "255"},
		{"1E", 16, "30"},
		{"101", 2, "5"},
		{"1010", 2, "10"},
		{"zz", 16, "zz"},
		{"2", 2, "2"},
		{"", 16, ""},
	}
	for _, c := range cases {
		if got := convertToDecimal(c.word, c.base); got != c.want {
			t.Errorf("convertToDecimal(%q, %d) = %q, want %q", c.word, c.base, got, c.want)
		}
	}
}

func TestCapitalize(t *testing.T) {
	cases := map[string]string{
		"bridge": "Bridge",
		"bRIDGE": "Bridge",
		"éLAN":   "Élan",
		"go":     "Go",
		"'hello": "'Hello",
		"":       "",
	}
	for in, want := range cases {
		if got := capitalize(in); got != want {
			t.Errorf("capitalize(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFixPunctuation(t *testing.T) {
	cases := map[string]string{
		"chat ,chien":   "chat, chien",
		"oui ... non":   "oui... non",
		"BAMM !!":       "BAMM!!",
		"what ?!next":   "what?! next",
		"no punctation": "no punctation",
	}
	for in, want := range cases {
		if got := fixPunctuation(in); got != want {
			t.Errorf("fixPunctuation(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFixQuotes(t *testing.T) {
	cases := []struct {
		in, want []string
	}{
		{[]string{"'", "awesome", "'"}, []string{"'awesome'"}},
		{[]string{"'", "a", "b", "'", "c"}, []string{"'a", "b'", "c"}},
		{[]string{"don't", "go"}, []string{"don't", "go"}},
		{[]string{"lonely", "'"}, []string{"lonely", "'"}},
	}
	for _, c := range cases {
		in := append([]string(nil), c.in...)
		if got := fixQuotes(in); !reflect.DeepEqual(got, c.want) {
			t.Errorf("fixQuotes(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFixArticles(t *testing.T) {
	cases := []struct {
		in, want []string
	}{
		{[]string{"a", "untold"}, []string{"an", "untold"}},
		{[]string{"A", "Hour"}, []string{"An", "Hour"}},
		{[]string{"a", "cat"}, []string{"a", "cat"}},
		{[]string{"a"}, []string{"a"}},
		{[]string{"a", "'apple'"}, []string{"an", "'apple'"}},
	}
	for _, c := range cases {
		in := append([]string(nil), c.in...)
		if got := fixArticles(in); !reflect.DeepEqual(got, c.want) {
			t.Errorf("fixArticles(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
