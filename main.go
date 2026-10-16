// Command go-reloaded lit un fichier texte, applique les règles d'édition
// du sujet (marqueurs, ponctuation, guillemets, articles) et écrit le
// résultat dans un fichier de sortie.
//
// Usage :
//
//	go run . <fichier_entree> <fichier_sortie>
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

const usage = "usage: go run . <fichier_entree> <fichier_sortie>"

// main lance le programme : en cas d'erreur, un message est écrit sur
// la sortie d'erreur et le programme se termine avec le code 1.
func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "erreur :", err)
		os.Exit(1)
	}
}

// run enchaîne les étapes : arguments → lecture → traitement → écriture.
// Elle renvoie une erreur au lieu de quitter, ce qui la rend testable.
func run(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("nombre d'arguments incorrect (%d au lieu de 2)\n%s", len(args), usage)
	}
	input, output := args[0], args[1]
	if input == "" || output == "" {
		return errors.New("un nom de fichier est vide\n" + usage)
	}

	data, err := os.ReadFile(input)
	if errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("fichier d'entrée introuvable : %s", input)
	}
	if err != nil {
		return fmt.Errorf("impossible de lire le fichier d'entrée : %w", err)
	}

	if err := os.WriteFile(output, []byte(Process(string(data))), 0644); err != nil {
		return fmt.Errorf("impossible d'écrire le fichier de sortie : %w", err)
	}
	return nil
}
