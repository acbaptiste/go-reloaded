package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRunEndToEnd lit un vrai fichier, le traite et vérifie le fichier écrit.
func TestRunEndToEnd(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "sample.txt")
	output := filepath.Join(dir, "result.txt")

	content := "it (cap) was the worst of times (up) , it was\nSimply add 42 (hex) and 10 (bin)\nbearing a untold story ,what do you think ?\n"
	want := "It was the worst of TIMES, it was\nSimply add 66 and 2\nbearing an untold story, what do you think?\n"

	if err := os.WriteFile(input, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{input, output}); err != nil {
		t.Fatalf("run a renvoyé une erreur : %v", err)
	}

	got, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("fichier de sortie illisible : %v", err)
	}
	if string(got) != want {
		t.Errorf("sortie incorrecte\n got  %q\n want %q", got, want)
	}
}

// TestRunErrors vérifie que chaque erreur est signalée par un message clair.
func TestRunErrors(t *testing.T) {
	dir := t.TempDir()
	existing := filepath.Join(dir, "in.txt")
	if err := os.WriteFile(existing, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name    string
		args    []string
		message string
	}{
		{"aucun argument", nil, "nombre d'arguments incorrect"},
		{"un seul argument", []string{existing}, "nombre d'arguments incorrect"},
		{"trop d'arguments", []string{existing, "out.txt", "extra"}, "nombre d'arguments incorrect"},
		{"nom vide", []string{"", "out.txt"}, "nom de fichier est vide"},
		{"fichier absent", []string{filepath.Join(dir, "nope.txt"), filepath.Join(dir, "out.txt")}, "introuvable"},
		{"dossier de sortie absent", []string{existing, filepath.Join(dir, "missing", "out.txt")}, "impossible d'écrire"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := run(c.args)
			if err == nil {
				t.Fatalf("run(%q) : erreur attendue, aucune reçue", c.args)
			}
			if !strings.Contains(err.Error(), c.message) {
				t.Errorf("message %q ne contient pas %q", err.Error(), c.message)
			}
		})
	}
}
