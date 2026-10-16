package main

import "testing"

// TestProcessSubjectExamples vérifie les exemples officiels du sujet et de
// l'audit : la sortie doit être identique caractère par caractère.
func TestProcessSubjectExamples(t *testing.T) {
	cases := []struct{ name, in, want string }{
		{
			"audit 1 : casse avec et sans nombre",
			"If I make you BREAKFAST IN BED (low, 3) just say thank you instead of: how (cap) did you get in my house (up, 2) ?",
			"If I make you breakfast in bed just say thank you instead of: How did you get in MY HOUSE?",
		},
		{
			"audit 2 : bin et hex",
			"I have to pack 101 (bin) outfits. Packed 1a (hex) just to be sure",
			"I have to pack 5 outfits. Packed 26 just to be sure",
		},
		{
			"audit 3 : ponctuation",
			"Don not be sad ,because sad backwards is das . And das not good",
			"Don not be sad, because sad backwards is das. And das not good",
		},
		{
			"audit 4 : cap n, guillemets et articles",
			"harold wilson (cap, 2) : ' I am a optimist ,but a optimist who carries a raincoat . '",
			"Harold Wilson: 'I am an optimist, but an optimist who carries a raincoat.'",
		},
		{
			"exemple de la page de garde",
			"it (cap) was the worst of times (up) , it was",
			"It was the worst of TIMES, it was",
		},
		{
			"hex et bin dans la même phrase",
			"Simply add 42 (hex) and 10 (bin) and you will see the result is 68.",
			"Simply add 66 and 2 and you will see the result is 68.",
		},
		{
			"groupes de ponctuation",
			"Punctuation tests are ... kinda boring ,what do you think ?",
			"Punctuation tests are... kinda boring, what do you think?",
		},
		{
			"article devant un mot en u",
			"bearing a untold story ,what do you think ?",
			"bearing an untold story, what do you think?",
		},
	}
	runProcessCases(t, cases)
}

// TestProcessRules vérifie chacune des dix règles avec les exemples du sujet.
func TestProcessRules(t *testing.T) {
	cases := []struct{ name, in, want string }{
		{"hex", "1E (hex) files were added", "30 files were added"},
		{"bin", "It has been 10 (bin) years", "It has been 2 years"},
		{"up", "Ready, set, go (up) !", "Ready, set, GO!"},
		{"low", "I should stop SHOUTING (low)", "I should stop shouting"},
		{"cap", "Welcome to the Brooklyn bridge (cap)", "Welcome to the Brooklyn Bridge"},
		{"up n", "This is so exciting (up, 2)", "This is SO EXCITING"},
		{"low n", "THESE WORDS ARE LOUD (low, 3)", "THESE words are loud"},
		{"cap n", "it was the age of foolishness (cap, 6)", "It Was The Age Of Foolishness"},
		{"ponctuation", "I was sitting over there ,and then BAMM !!", "I was sitting over there, and then BAMM!!"},
		{"groupe ...", "I was thinking ... You were right", "I was thinking... You were right"},
		{"groupe !?", "Really !? Yes", "Really!? Yes"},
		{"deux-points et point-virgule", "note : this ;that", "note: this; that"},
		{"guillemets un mot", "I am exactly how they describe me: ' awesome '", "I am exactly how they describe me: 'awesome'"},
		{"guillemets plusieurs mots", "As Elton John said: ' I am the most well-known homosexual in the world '", "As Elton John said: 'I am the most well-known homosexual in the world'"},
		{"a → An", "There it was. A amazing rock!", "There it was. An amazing rock!"},
		{"a → an devant h", "a hour ago", "an hour ago"},
		{"a reste a", "a cat and a dog", "a cat and a dog"},
	}
	runProcessCases(t, cases)
}

// TestProcessEdgeCases couvre les cas limites listés dans l'analyse du sujet.
func TestProcessEdgeCases(t *testing.T) {
	cases := []struct{ name, in, want string }{
		{"chaîne vide", "", ""},
		{"marqueur seul", "(up)", ""},
		{"marqueur en début de ligne", "(up) hello", "hello"},
		{"marqueur collé au mot", "word(up) here", "WORD here"},
		{"marqueur sans espace", "so exciting (up,2)", "SO EXCITING"},
		{"marqueur avec espaces internes", "so exciting ( up , 2 )", "SO EXCITING"},
		{"n trop grand", "go now (up, 9)", "GO NOW"},
		{"n énorme (dépassement int)", "go now (up, 99999999999999999999999)", "GO NOW"},
		{"n nul", "keep THIS (low, 0)", "keep THIS"},
		{"n négatif non reconnu", "hello (up, -1)", "hello (up, -1)"},
		{"marqueur inconnu laissé tel quel", "hello (rev)", "hello (rev)"},
		{"marqueurs successifs", "go (up) (low)", "go"},
		{"ponctuation ignorée par le comptage", "a, b (cap, 2)", "A, B"},
		{"hex invalide laissé intact", "zz (hex) end", "zz end"},
		{"bin invalide laissé intact", "102 (bin) end", "102 end"},
		{"hex avec n", "a 1f 10 (hex, 2)", "a 31 16"},
		{"cap sur mot en majuscules", "bRIDGE (cap)", "Bridge"},
		{"accents", "élan (cap) déjà (up)", "Élan DÉJÀ"},
		{"apostrophe dans un mot", "don't stop (up, 2)", "DON'T STOP"},
		{"apostrophe non prise pour guillemet", "it's ' fine '", "it's 'fine'"},
		{"ponctuation avant guillemet fermant", "' hello world ! '", "'hello world!'"},
		{"article devant guillemet", "a ' amazing ' thing", "an 'amazing' thing"},
		{"article après conversion de casse", "a APPLE (low)", "an apple"},
		{"espaces multiples", "hello     world   ,  again", "hello world, again"},
		{"retours à la ligne conservés", "first line (up)\nsecond ,line", "first LINE\nsecond, line"},
		{"lignes vides conservées", "a\n\nb", "a\n\nb"},
		{"fins de ligne Windows", "hello (up)\r\nworld", "HELLO\nworld"},
	}
	runProcessCases(t, cases)
}

func runProcessCases(t *testing.T, cases []struct{ name, in, want string }) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Process(c.in); got != c.want {
				t.Errorf("Process(%q)\n got  %q\n want %q", c.in, got, c.want)
			}
		})
	}
}
