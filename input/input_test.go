package input

import (
	"reflect"
	"strings"
	"testing"
)

var badInputs = []struct {
	in        string
	errorType string
}{
	{"", "BLANK INPUT"},
	{",", "BLANK GLYPH NAME"},
	{" ", "BLANK INPUT"},
	{"\r", "BLANK INPUT"},
	{"\n", "BLANK INPUT"},
	{"\t", "BLANK INPUT"},
	{"\r\n\r\n\r\n", "BLANK INPUT"},
	{"\r\n\r\n\r\n", "BLANK INPUT"},
	{"\r\r\r\r\r\r", "BLANK INPUT"},
	{"\n\n\n\n\n\n", "BLANK INPUT"},
	{"\r\t\n\r\t\n \t\t  ", "BLANK INPUT"},
	{"1,2,3,4,5,6,7,8,9,10", "TOO MANY GLYPHS"},
	{"xm,xm,xm,xm,xm,xm,xm,xm,xm,xm", "TOO MANY GLYPHS"},
	{"This string is 201 characters long. The quick brown fox jumped over the lazy dog. The quick brown fox jumped over the lazy dog. The quick brown fox jumped over the lazy dog. The quick brown fox jumped.", "INPUT TOO LONG"},
	{",,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,", "TOO MANY GLYPHS"},
	{"你好, 世界", "UNKNOWN GLYPH NAME"},         // Test Unicode handling. Hello, world! / Nǐ hǎo, shìjiè!
	{"Nǐ hǎo, shìjiè", "UNKNOWN GLYPH NAME"}, //Hello, world! / Nǐ hǎo, shìjiè!
	{"你好", "UNKNOWN GLYPH NAME"},             // Hello!
	{"\xc3\x28", "INVALID UTF-8"},            // Invalid UTF-8 http://stackoverflow.com/questions/1301402/example-invalid-utf8-string
	{"\xff\xfe\xfd", "INVALID UTF-8"},        // Canonical invalid UTF-8 from https://golang.org/pkg/unicode/utf8/#example_Valid
}

func TestBadInputs(t *testing.T) {
	for _, tt := range badInputs {
		_, _, err := ProcessString(tt.in)
		if err == nil {
			t.Errorf("Bad input didn't get rejected! Expected error %s for input `%s`.", tt.errorType, tt.in)
			continue
		}
		if strings.HasPrefix(err.Error(), tt.errorType) {
			//t.Logf("Bad input successfully rejected because %s: `%s`", tt.errorType, tt.in)
			continue
		}
		t.Errorf("Bad input rejected with wrong error type! Got error %s. Expected error %s for input `%s`.", err.Error(), tt.errorType, tt.in)
	}
}

var goodInputs = []struct {
	in         string
	glyphNames []string
	edgeLists  []string
}{
	{"XM", []string{"XM"}, []string{"67697a898a"}}, // Simple case
	{"xm", []string{"XM"}, []string{"67697a898a"}}, // Simple case
	{"xM", []string{"XM"}, []string{"67697a898a"}}, // Simple case
	{"Xm", []string{"XM"}, []string{"67697a898a"}}, // Simple case
	{
		"Forget, Old, See, New",
		[]string{"Forget", "Old", "See", "New"},
		[]string{"48", "5989", "09", "2767"},
	},
	{ // No whitespace
		"Forget,Old,See,New",
		[]string{"Forget", "Old", "See", "New"},
		[]string{"48", "5989", "09", "2767"},
	},
	{ // All lowercase
		"forget, old, see, new",
		[]string{"Forget", "Old", "See", "New"},
		[]string{"48", "5989", "09", "2767"},
	},
	{ // All uppercase
		"FORGET, OLD, SEE, NEW",
		[]string{"Forget", "Old", "See", "New"},
		[]string{"48", "5989", "09", "2767"},
	},
	{ // weirdly mixed case
		"fORGEt, oLd, sEe, nEw",
		[]string{"Forget", "Old", "See", "New"},
		[]string{"48", "5989", "09", "2767"},
	},
	{ // Whitespace abuse
		" \t\r\nForget \t\r\n, \t\r\nOld \t\r\n, \t\r\nSee \t\r\n, \t\r\nNew \t\r\n",
		[]string{"Forget", "Old", "See", "New"},
		[]string{"48", "5989", "09", "2767"},
	},
	{ // Test the ones that have spaces in the glyph name
		"Open All, Clear All, Discover, Truth",
		[]string{"Open All", "Clear All", "Discover", "Truth"},
		[]string{"010512233437384578", "01050a1223343a45", "122334", "676a7a898a9a"},
	},
	{
		"N'zeer, Openall, Clearall, XM",
		[]string{"N'zeer", "Openall", "Clearall", "XM"},
		[]string{"06090a3a6a9a", "010512233437384578", "01050a1223343a45", "67697a898a"},
	},
	// Wrong separators (spaces instead of commas)
	// Wrong separators (semicolons instead of commas)
	// Accidental plurals (portals vs portal)

}

func TestGoodInputs(t *testing.T) {
	for _, tt := range goodInputs {
		g, e, err := ProcessString(tt.in)
		if err != nil {
			t.Errorf("Input produced an unexpected error.\nInput: `%s`\nError: `%s`", tt.in, err.Error())
			continue
		}
		if !(reflect.DeepEqual(g, tt.glyphNames)) || !(reflect.DeepEqual(e, tt.edgeLists)) {
			t.Errorf("Bad return values.\nInput: `%s`\nGot: `%s`, `%s`\nExpected: `%s`, `%s`",
				tt.in, g, e, tt.glyphNames, tt.edgeLists)
			continue
		}

		//t.Logf("Case `%s` -> %s, %s passed!", tt.in, tt.glyphNames, tt.edgeLists)

	}
}

var marginalInputs = []struct {
	in        string
	errorType string
}{
	{
		"Help Enlightened Capture All Portal",
		"USE COMMAS, NOT SPACES",
	},
	{
		"Help;Enlightened;Capture;All;Portal",
		"USE COMMAS, NOT SEMICOLONS",
	},
	{
		"Help-Enlightened-Capture-All-Portal",
		"USE COMMAS, NOT DASHES",
	},
}

func TestMarginalInputs(t *testing.T) {
	for _, tt := range marginalInputs {
		_, _, err := ProcessString(tt.in)
		if err == nil {
			t.Error("Failed to reject input", tt.in, "Expected error", tt.errorType)
		}
		if !strings.HasPrefix(err.Error(), tt.errorType) {
			t.Error("Failed with wrong error", err.Error(), "Expected error", tt.errorType)
		}
	}
}

var futureInputs = []struct {
	in        string
	errorType string
}{
	{
		"PORTALS",
		"UNKNOWN GLYPH NAME: Did you mean Portal?",
	},
}

// Expect this to fail until Levenshtein distance corrections are implemented.
// The program should suggest the closest spelling to you.
func TestFutureInputs(t *testing.T) {
	for _, tt := range futureInputs {
		_, _, err := ProcessString(tt.in)
		if err == nil {
			t.Error("Failed to reject input", tt.in, "Expected error", tt.errorType)
		}
		if !strings.HasPrefix(err.Error(), tt.errorType) {
			t.Error("Failed with wrong error", err.Error(), "Expected error", tt.errorType)
		}
	}
}
