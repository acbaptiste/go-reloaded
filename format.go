package main

import (
	"regexp"
	"strings"
)

var (
	// spaceBeforePunctuation : des blancs suivis d'un groupe de signes.
	spaceBeforePunctuation = regexp.MustCompile(`\s+([.,!?:;]+)`)
	// missingSpaceAfter : un groupe de signes collé à autre chose qu'un blanc.
	missingSpaceAfter = regexp.MustCompile(`([.,!?:;]+)([^\s.,!?:;])`)
)

// fixPunctuation colle . , ! ? : ; au mot précédent et met un espace avant
// le mot suivant. Le "+" des motifs garde ensemble les groupes "..." et "!?".
// Go (RE2) n'a pas de lookahead : on procède donc en deux passes.
func fixPunctuation(text string) string {
	text = spaceBeforePunctuation.ReplaceAllString(text, "$1")
	return missingSpaceAfter.ReplaceAllString(text, "$1 $2")
}

// fixQuotes colle chaque paire de ' isolés aux mots qu'elle entoure :
// le premier ' se colle au mot suivant, le second au mot précédent.
// Seuls les tokens égaux à "'" sont concernés : l'apostrophe de "don't"
// n'est jamais prise pour un guillemet.
func fixQuotes(tokens []string) []string {
	out := make([]string, 0, len(tokens))
	open := false
	for i, token := range tokens {
		switch {
		case token != "'":
			out = append(out, token)
		case !open && i+1 < len(tokens):
			// Guillemet ouvrant : on le colle au token suivant.
			tokens[i+1] = "'" + tokens[i+1]
			open = true
		case open && len(out) > 0:
			// Guillemet fermant : on le colle au dernier token écrit.
			out[len(out)-1] += "'"
			open = false
		default:
			// Guillemet seul sans partenaire possible : on le garde tel quel.
			out = append(out, token)
		}
	}
	return out
}

// fixArticles remplace "a" par "an" (et "A" par "An") quand le mot suivant
// commence par une voyelle (a, e, i, o, u) ou un h.
func fixArticles(tokens []string) []string {
	for i := 0; i < len(tokens)-1; i++ {
		if (tokens[i] == "a" || tokens[i] == "A") && startsWithVowelOrH(tokens[i+1]) {
			tokens[i] += "n"
		}
	}
	return tokens
}

// startsWithVowelOrH regarde la première lettre du mot en ignorant un
// éventuel guillemet ouvrant ('amazing') ; la casse n'a pas d'importance.
func startsWithVowelOrH(word string) bool {
	word = strings.TrimLeft(word, `'"`)
	return word != "" && strings.ContainsRune("aeiouhAEIOUH", rune(word[0]))
}
