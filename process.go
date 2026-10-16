package main

import "strings"

// Process applique toutes les règles du sujet à un texte complet.
// Le texte est traité ligne par ligne pour conserver les retours à la ligne
// (strings.Fields, utilisé plus loin, avalerait les "\n").
func Process(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = processLine(strings.TrimSuffix(line, "\r"))
	}
	return strings.Join(lines, "\n")
}

// processLine fait traverser une ligne au pipeline, dans cet ordre :
//
//  1. marqueurs   : (hex) (bin) (up) (low) (cap) et leurs variantes (cmd, n)
//  2. ponctuation : . , ! ? : ; collés au mot précédent
//  3. guillemets  : ' ... ' collés aux mots qu'ils entourent
//  4. articles    : a → an devant une voyelle ou un h
//
// Les marqueurs passent en premier car ils doivent disparaître avant la mise
// en forme ; les articles en dernier car ils regardent le mot suivant tel
// qu'il sera écrit.
func processLine(line string) string {
	tokens := applyMarkers(line)
	text := fixPunctuation(strings.Join(tokens, " "))
	tokens = fixQuotes(strings.Fields(text))
	tokens = fixArticles(tokens)
	return strings.Join(tokens, " ")
}
