package input

import (
	"errors"
	"fmt"
	"github.com/LiaungYip/glyphs"
	"strings"
	"unicode/utf8"
)

// Do basic input sanity checks.
// 1. Input is less than 200 bytes (prevent DoS attack by huge inputs)
// 2. Input is valid UTF-8
// 3. Input isn't blank
func sanityCheck(xptr *string) error {
	if len(*xptr) > 200 {
		return errors.New("INPUT TOO LONG: You sent me too much text! I can't handle that.")
	}

	if !utf8.ValidString(*xptr) {
		return errors.New("INVALID UTF-8")
	}

	if len(strings.TrimSpace(*xptr)) == 0 {
		return errors.New("BLANK INPUT: You sent me an empty message! I can't do anything with that.")
	}
	return nil
}

// Split input string into words.
// Sanity check that there are between 1 and 9 words.
func splitWords(xptr *string, sep string) ([]string, error) {
	words := strings.Split(*xptr, sep)
	numWords := len(words)
	for n, w := range words {
		trimmed := strings.TrimSpace(w)
		if trimmed == "" {
			return nil, errors.New("BLANK GLYPH NAME: Did you have too many commas?")
		}
		words[n] = trimmed
	}
	if numWords > 9 {
		s := fmt.Sprintf("TOO MANY GLYPHS: You can have a maximum of 9 glyphs. (You sent %d glyphs.)", len(words))
		return nil, errors.New(s)
	}
	if numWords < 1 {
		return nil, errors.New("BLANK INPUT: You sent me an empty message! I can't do anything with that.")
	}

	return words, nil
}

func lookupWords(words []string) ([]string, []string, []error) {
	numWords := len(words)
	glyphNames := make([]string, numWords)
	edgeLists := make([]string, numWords)
	errs := make([]error, numWords)

	for n, w := range words {
		glyphNames[n], edgeLists[n], errs[n] = glyphs.Lookup(w)
	}

	return glyphNames, edgeLists, errs
}

func concatenateErrors(errs []error) error {
	error_happened := false
	var err_strings []string
	for _, e := range errs {
		if e != nil {
			err_strings = append(err_strings, e.Error())
			error_happened = true
		}
	}

	if error_happened {
		err_string := strings.Join(err_strings, "\n")
		return errors.New(err_string)
	} else {
		return nil
	}
}

var badSeparators = []struct {
	sep  string
	name string
}{
	{" ", "SPACES"},
	{";", "SEMICOLONS"},
	{"-", "DASHES"},
}

// ProcessInputString accepts a comma-separated string of glyph names.
// On success, it returns arrays of glyph names, and edge lists.
// On failure, it returns a descriptive error message to give back to the user.
func ProcessString(x string) ([]string, []string, error) {
	err := sanityCheck(&x)
	if err != nil {
		return nil, nil, err
	}

	words, err := splitWords(&x, ",")
	if err != nil {
		return nil, nil, err
	}

	glyphNames, edgeLists, errs := lookupWords(words)

	err = concatenateErrors(errs)
	if err == nil { // Success
		return glyphNames, edgeLists, nil
	}

	// ---------------- Check for common user input errors. ----------------

	// Trial using ` `, `;`, and `-` as separators.
	// Tell the user if the input parses correctly using one of these separators.
	for _, s := range badSeparators {
		words, err := splitWords(&x, s.sep)
		if err != nil {
			continue
		}
		_, _, errs := lookupWords(words)
		err2 := concatenateErrors(errs)
		if err2 == nil {
			t := "USE COMMAS, NOT " + s.name + ". Example: `Clear All, Open All, Discover, Truth`."
			return nil, nil, errors.New(t)
		}
	}

	for _, w := range words {
		candidate := glyphs.Spellcheck(w)
		if candidate != "" && !strings.EqualFold(candidate, w) {
			t := fmt.Sprintf("UNKNOWN GLYPH NAME: %s. Did you mean %s?", w, candidate)
			return nil, nil, errors.New(t)
		}
	}

	return nil, nil, err
}
