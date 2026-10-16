package main

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// markerPattern reconnaît un marqueur, avec ou sans nombre et avec des
// espaces optionnels : (up)  (cap, 6)  (low,2)  ( hex ).
//
//	groupe 1 : la commande   groupe 2 : le nombre (vide s'il est absent)
var markerPattern = regexp.MustCompile(`\(\s*(hex|bin|up|low|cap)\s*(?:,\s*(\d+)\s*)?\)`)

// applyMarkers découpe une ligne en mots et applique chaque marqueur aux
// mots qui le précèdent ; les marqueurs eux-mêmes disparaissent.
//
// La regex coupe la ligne autour des marqueurs : Split donne les morceaux
// de texte entre eux, FindAllStringSubmatch donne les marqueurs. Il y a
// toujours un morceau de texte de plus que de marqueurs :
//
//	"so exciting (up, 2) now" → texts   = ["so exciting ", " now"]
//	                            markers = [["(up, 2)", "up", "2"]]
func applyMarkers(line string) []string {
	texts := markerPattern.Split(line, -1)
	markers := markerPattern.FindAllStringSubmatch(line, -1)

	words := strings.Fields(texts[0])
	for i, m := range markers {
		applyToPreviousWords(words, m[1], markerCount(m[2]))
		words = append(words, strings.Fields(texts[i+1])...)
	}
	return words
}

// markerCount lit le nombre d'un marqueur : 1 s'il est absent, le plus
// grand int possible s'il est trop grand (on modifie alors tous les mots).
func markerCount(number string) int {
	if number == "" {
		return 1
	}
	n, err := strconv.Atoi(number)
	if err != nil {
		return math.MaxInt
	}
	return n
}

// applyToPreviousWords remonte dans words et transforme les count derniers
// mots. Les tokens de pure ponctuation (",", "!", "'") ne comptent pas.
// La boucle s'arrête d'elle-même si count dépasse le nombre de mots.
func applyToPreviousWords(words []string, command string, count int) {
	for i := len(words) - 1; i >= 0 && count > 0; i-- {
		if isWord(words[i]) {
			words[i] = transformWord(command, words[i])
			count--
		}
	}
}

// transformWord applique une commande à un seul mot.
func transformWord(command, word string) string {
	switch command {
	case "hex":
		return convertToDecimal(word, 16)
	case "bin":
		return convertToDecimal(word, 2)
	case "up":
		return strings.ToUpper(word)
	case "low":
		return strings.ToLower(word)
	case "cap":
		return capitalize(word)
	}
	return word
}

// convertToDecimal convertit un nombre écrit dans la base donnée en décimal.
// Si le mot n'est pas un nombre valide, il est laissé intact (pas de panic).
func convertToDecimal(word string, base int) string {
	n, err := strconv.ParseInt(word, base, 64)
	if err != nil {
		return word
	}
	return strconv.FormatInt(n, 10)
}

// capitalize met la première lettre en majuscule et le reste en minuscules.
// On travaille sur des runes pour ne pas casser les accents (élan → Élan),
// et on saute un éventuel guillemet ou signe en tête ('hello → 'Hello).
func capitalize(word string) string {
	runes := []rune(strings.ToLower(word))
	for i, r := range runes {
		if isLetterOrDigit(r) {
			runes[i] = unicode.ToUpper(r)
			break
		}
	}
	return string(runes)
}

// isWord indique si un token contient au moins une lettre ou un chiffre.
func isWord(token string) bool {
	return strings.IndexFunc(token, isLetterOrDigit) >= 0
}

func isLetterOrDigit(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}
