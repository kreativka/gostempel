package gostempel

import (
	"testing"
)

var stemTest = []struct {
	in  string
	out string
}{
	{"aachen", "aachneć"},
	{"aarem", "aar"},
	{"aasen", "aąć"},
	{"abecadłarzów", "abecadłarz"},
	{"bigosującemu", "bigosujący"},
	{"certyfikujesz", "certyfikować"},
	{"czterokilometrowego", "czterokilometrowy"},
	{"domińskimi", "domiński"},
	{"emilciny", "emeć"},
	{"fotyczność", "fotycznoć"},
	{"hamulcowych", "hamulcowy"},
	{"kurytybsko", "kurytybsko"},
	{"lupka", "lupeky"},
	{"międzykomórkowość", "międzykomórkowość"},
	{"neohellenizm", "neohellenizm"},
	{"niegnębicielska", "niegnębicielski"},
	{"niekolebkowi", "niekok"},
	{"nienyskości", "nyskość"},
	{"niewirginijski", "niewirginijski"},
	{"niezgodność", "niezgodnć"},
	{"niezgodny", "nieza"},
	{"nieodprzysięgniętemu", "nieodprzysięgniętemć"},
	{"piechurskiego", "piechurski"},
	{"piędzią", "piędź"},
	{"powstawiajmy", "powstawiać"},
	{"synchrotronową", "synchrotronowy"},
	{"szpakowiczów", "szpakowicz"},
	{"wielobrzmiącej", "wielobrzmiący"},
	{"ziemaka", "ziemak"},
	{"ekstensjonalny", "ekstensjonalnć"},
	{"mccullersowie", "mccullers"},
	{"nabierający", "nabierający"},
	{"niepopoborowego", "niepopoborowy"},
	{"proceduralne", "procedurcny"},
	{"sześciotygodniowego", "sześciotygodniowy"},
	{"wstawowość", "wstawowość"},
	{"zakulisowych", "zakulisowy"},
	{"barlinkiem", "barlinek"},
	{"niehamletowaniem", "niehamletować"},
	{"nienaczynanej", "nienaczynać"},
	{"nieprzyłącznego", "nieprzyłąóąy"},
	{"nieprzyłączny", "nieprzyłączn"},
	{"oblesioność", "oblesioć"},
	{"piekielnikami", "piekielnik"},
	{"przekłamanymi", "przekłamać"},
	{"salwickimi", "salwicki"},
	{"walkerów", "walker"},
	{"ważniejszemu", "oważny"},
	{"wyprzedawać", "wyprzedawać"},
	{"żyznościowym", "żyznościowy"},
	// For testing set minTokenLength to 0 and
	// uncomment below tokens
	//{"żzw", "żzwa"},
}

var f = "./stemmer_20000.tbl"
var stemmer, _ = LoadStemmer(f)

func TestStem(t *testing.T) {
	for _, tt := range stemTest {
		s := Stem(stemmer, []rune(tt.in))
		if string(s) != tt.out {
			t.Errorf("stem(%q) => %q, want %q", tt.in, s, tt.out)
		}
	}
}

func BenchmarkLoadStemmer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stemmer, _ = LoadStemmer(f)
	}
}

func benchmarkStem(s []rune, b *testing.B) {
	for i := 0; i < b.N; i++ {
		Stem(stemmer, s)
	}
}

func BenchmarkStemPowyciagajacy(b *testing.B)  { benchmarkStem([]rune("powyciągający"), b) }
func BenchmarkStemAchen(b *testing.B)          { benchmarkStem([]rune("aachen"), b) }
func BenchmarkStemCiupciajacy(b *testing.B)    { benchmarkStem([]rune("ciupciający"), b) }
func BenchmarkStemZmiim(b *testing.B)          { benchmarkStem([]rune("żmiim"), b) }
func BenchmarkStemKupilem(b *testing.B)        { benchmarkStem([]rune("kupiłem"), b) }
func BenchmarkStemUderzajac(b *testing.B)      { benchmarkStem([]rune("uderzając"), b) }
func BenchmarkStemNiedokrwistosc(b *testing.B) { benchmarkStem([]rune("niedokrwistość"), b) }
