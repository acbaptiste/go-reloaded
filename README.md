# go-reloaded

Outil en ligne de commande écrit en Go. Il lit un fichier texte, corrige le texte selon des règles d'édition et écrit le résultat dans un second fichier.

```
it (cap) was the worst of times (up) , it was      →   It was the worst of TIMES, it was
Simply add 42 (hex) and 10 (bin)                    →   Simply add 66 and 2
bearing a untold story ,what do you think ?         →   bearing an untold story, what do you think?
```

Le projet n'utilise que la bibliothèque standard de Go.

## Prérequis

- Go 1.22 ou plus récent (`go version` pour vérifier)

## Installation

```bash
git clone <url-du-depot> go-reloaded
cd go-reloaded
```

## Compilation

```bash
go build -o go-reloaded .
```

## Utilisation

```bash
go run . <fichier_entree> <fichier_sortie>
# ou, après compilation :
./go-reloaded <fichier_entree> <fichier_sortie>
```

- Le premier argument est le fichier à corriger. Il doit exister.
- Le second argument est le fichier où écrire le résultat. Il est créé (droits `0644`) ou écrasé.

En cas d'erreur, le programme écrit un message sur la sortie d'erreur et se termine avec le code `1`. Les cas couverts sont un nombre d'arguments incorrect, un fichier d'entrée introuvable ou illisible, et un fichier de sortie impossible à écrire.

```bash
$ go run . sample.txt
erreur : nombre d'arguments incorrect (1 au lieu de 2)
usage: go run . <fichier_entree> <fichier_sortie>
```

## Fonctionnalités

| Règle | Entrée | Sortie |
|---|---|---|
| `(hex)` : hexadécimal vers décimal | `1E (hex) files were added` | `30 files were added` |
| `(bin)` : binaire vers décimal | `It has been 10 (bin) years` | `It has been 2 years` |
| `(up)` : majuscules | `Ready, set, go (up) !` | `Ready, set, GO!` |
| `(low)` : minuscules | `I should stop SHOUTING (low)` | `I should stop shouting` |
| `(cap)` : première lettre en majuscule | `Welcome to the Brooklyn bridge (cap)` | `Welcome to the Brooklyn Bridge` |
| `(up, n)`, `(low, n)`, `(cap, n)` : les n mots précédents | `This is so exciting (up, 2)` | `This is SO EXCITING` |
| Ponctuation `. , ! ? : ;` collée au mot précédent | `over there ,and then BAMM !!` | `over there, and then BAMM!!` |
| Groupes `...` et `!?` conservés | `I was thinking ... You were right` | `I was thinking... You were right` |
| Paires de `'` collées aux mots entourés | `me: ' awesome '` | `me: 'awesome'` |
| `a` devient `an` devant a, e, i, o, u ou h | `A amazing rock!` | `An amazing rock!` |

## Exemples

Le dossier [`examples/`](examples) contient un fichier d'entrée et la sortie attendue :

```bash
go run . examples/sample.txt result.txt
diff result.txt examples/result.txt   # aucune différence attendue
```

## Tests

```bash
go test ./...          # lancer tous les tests
go test -v ./...       # avec le détail de chaque cas
go test -cover ./...   # avec le taux de couverture
go vet ./...           # analyse statique
gofmt -l .             # doit n'afficher aucun fichier
```

Les tests sont table-driven : chaque cas est une ligne avec une entrée et une sortie attendue. Ajouter un cas revient à ajouter une ligne.

| Fichier | Ce qu'il vérifie |
|---|---|
| `process_test.go` | Les exemples officiels du sujet et de l'audit, chacune des dix règles, et les cas limites (voir plus bas) |
| `helpers_test.go` | Chaque fonction isolément : application des marqueurs, lecture du nombre, conversions, capitalisation, ponctuation, guillemets, articles |
| `main_test.go` | Le programme de bout en bout avec de vrais fichiers, et chaque cas d'erreur (arguments, fichier absent, écriture impossible) |

## Organisation du code

```
go-reloaded/
├── main.go           entrées / sorties : run vérifie les arguments, lit, traite, écrit
├── process.go        Process : découpage en lignes et pipeline de traitement
├── markers.go        étape 1 : marqueurs (hex) (bin) (up) (low) (cap) et (cmd, n)
├── format.go         étapes 2 à 4 : ponctuation, guillemets, articles
├── *_test.go         tests
└── examples/         fichier d'exemple et sortie attendue
```

`main` ne fait qu'appeler `run`, qui enchaîne les étapes :

```
arguments → vérification → os.ReadFile → Process → os.WriteFile → fichier de sortie
```

Toute la logique de texte est dans `Process`, une fonction pure : elle prend une chaîne et renvoie une chaîne. On peut donc la tester sans créer de fichier.

## Choix techniques

**Traitement ligne par ligne.** `strings.Fields` supprime les retours à la ligne. Le texte est donc découpé en lignes avant le traitement, puis les lignes sont recollées avec `"\n"`. La structure du fichier est ainsi conservée.

**Ordre du pipeline.** Chaque ligne passe par quatre étapes, dans cet ordre :

1. **Marqueurs.** Ils doivent disparaître avant la mise en forme. Dans `go (up) !`, le `!` doit se coller à `GO` et non à `(up)`.
2. **Ponctuation.** Une fois les marqueurs retirés, il ne reste que des mots et des signes.
3. **Guillemets.** Le signe est collé avant le guillemet fermant : `raincoat . '` devient `raincoat. '` puis `raincoat.'`.
4. **Articles.** La règle regarde le mot suivant tel qu'il sera écrit, après les conversions et les changements de casse.

**Marqueurs avec nombre.** Un découpage naïf coupe `(up, 2)` en deux tokens, `"(up,"` et `"2)"`. Une expression régulière repère donc les marqueurs directement dans la ligne : `Split` donne les morceaux de texte entre les marqueurs, `FindAllStringSubmatch` donne les marqueurs eux-mêmes. Seuls les morceaux de texte sont découpés en mots. La regex accepte aussi `(up,2)`, `( up , 2 )` et un marqueur collé au mot (`word(up)`).

**Parcours arrière.** Quand on rencontre un marqueur, on remonte dans les mots déjà écrits pour modifier les n derniers. Les tokens de pure ponctuation (`,`, `!`) ne comptent pas comme des mots.

**Ponctuation en deux passes.** Le moteur d'expressions régulières de Go (RE2) n'a pas de lookahead. La première passe retire l'espace avant les signes, la seconde ajoute l'espace après. Le `+` du motif garde les groupes `...` et `!?` ensemble.

**Guillemets.** Seuls les `'` isolés sont des guillemets, donc l'apostrophe de `don't` n'est jamais touchée. Un booléen indique si le prochain `'` ouvre ou ferme.

**Unicode.** `(cap)` travaille sur des runes et non sur des octets, ce qui donne `élan` → `Élan`.

## Cas limites et comportements choisis

| Situation | Exemple | Comportement |
|---|---|---|
| Marqueur en début de ligne | `(up) hello` | Le marqueur est supprimé, rien n'est modifié |
| n plus grand que le nombre de mots | `go now (up, 9)` | Tous les mots disponibles sont modifiés : `GO NOW` |
| n égal à 0 | `keep THIS (low, 0)` | Aucun mot modifié |
| n négatif ou non numérique | `hello (up, -1)` | Ce n'est pas un marqueur, le texte est laissé tel quel |
| Marqueur inconnu | `hello (rev)` | Laissé tel quel |
| Plusieurs marqueurs de suite | `go (up) (low)` | Appliqués dans l'ordre : `go` |
| Nombre invalide pour `(hex)` / `(bin)` | `zz (hex)` | Le mot est laissé intact, sans plantage |
| Fins de ligne Windows | `\r\n` | Converties en `\n` |

## Limites connues

- Un nombre décimal ou avec séparateur de milliers est coupé par la règle de ponctuation : `3.14` devient `3. 14`. Le sujet ne demande pas de traiter ce cas.
- Les espaces en début et en fin de ligne, ainsi que les espaces multiples entre les mots, sont réduits à un seul espace.
- Une paire de guillemets ne doit pas être répartie sur deux lignes. L'état ouvert/fermé est remis à zéro à chaque ligne.
