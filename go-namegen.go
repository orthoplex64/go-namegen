package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/orthoplex64/go-namegen/distr"
)

var (
	flagN                 = flag.Int("n", 10, "number of names to generate")
	flagDistrVowels       = flag.String("voweldistr", "english", "commma-separated list of colon-separated pairs of vowels and their weights")
	flagDistrConsonants   = flag.String("consdistr", "english", "commma-separated list of colon-separated pairs of consonants and their weights")
	flagDistrSyllables    = flag.String("syldistr", "cv:2,cvc,vc:2", "commma-separated list of colon-separated pairs of types of syllables and their weights")
	flagDistrNumSyllables = flag.String("numsyldistr", "2,3:3,4", "commma-separated list of colon-separated pairs of numbers of syllables and their weights")

	distrVowels      *distr.StrDistr
	distrConsonants  *distr.StrDistr
	distrLetters     *distr.StrDistr
	distrSyllables   *distr.StrDistr
	distrNumSyllales *distr.StrDistr

	strIntCache map[string]int
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintln(os.Stderr, `"english" and "uniform" are also valid for -voweldistr and -consdistr.`)
		fmt.Fprintln(os.Stderr, `Syllable types accepted by -syldistr are strings of 'v' (vowel), 'c' (consonant), and/or 'l' (letter).`)
		flag.PrintDefaults()
	}
}

func parseMapFlag(s string) map[string]float64 {
	words := strings.Split(strings.TrimSpace(s), ",")
	m := make(map[string]float64, len(words))
	for _, word := range words {
		pair := strings.Split(strings.TrimSpace(word), ":")
		switch len(pair) {
		case 1:
			m[pair[0]]++
		case 2:
			val, err := strconv.ParseFloat(pair[1], 64)
			if err != nil {
				log.Fatalf("parseMapFlag: error parsing %q: %v\n", pair[1], err)
			}
			m[pair[0]] += val
		default:
			log.Fatalf("parseMapFlag: too many colons in %q\n", word)
		}
	}
	return m
}

func parseFlags() {
	flag.Parse()
	if !flag.Parsed() {
		flag.PrintDefaults()
		os.Exit(0)
	}
	// vowel distribution
	switch *flagDistrVowels {
	case "uniform":
		distrVowels = uniformVowels
	case "english":
		distrVowels = englishVowels
	default:
		m := parseMapFlag(*flagDistrVowels)
		distrVowels = distr.NewStrDistr()
		for s, p := range m {
			for _, c := range s {
				distrVowels.Add(string(c), p)
			}
		}
	}
	// consonant distribution
	switch *flagDistrConsonants {
	case "uniform":
		distrConsonants = uniformConsonants
	case "english":
		distrConsonants = englishConsonants
	default:
		m := parseMapFlag(*flagDistrConsonants)
		distrConsonants = distr.NewStrDistr()
		for s, p := range m {
			for _, c := range s {
				distrConsonants.Add(string(c), p)
			}
		}
	}
	// letter distribution
	distrLetters = distr.NewStrDistr()
	for _, d := range []*distr.StrDistr{distrVowels, distrConsonants} {
		for _, s := range d.Strings() {
			distrLetters.Add(s, d.Weight(s))
		}
	}
	// syllable distribution
	distrSyllables = distr.NewStrDistr()
	for s, p := range parseMapFlag(*flagDistrSyllables) {
		distrSyllables.Add(s, p)
	}
	// syllable count distribution
	distrNumSyllales = distr.NewStrDistr()
	m := parseMapFlag(*flagDistrNumSyllables)
	strIntCache = make(map[string]int, len(m))
	for s, p := range m {
		distrNumSyllales.Add(s, p)
		val, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalf("parseFlags: error parsing %q: %v\n", s, err)
		}
		strIntCache[s] = val
	}
}

func main() {
	parseFlags()
	for i := 0; i < *flagN; i++ {
		buf := new(bytes.Buffer)
		numSyls := strIntCache[distrNumSyllales.Pick()]
		for syl := 0; syl < numSyls; syl++ {
			sylType := distrSyllables.Pick()
			for _, c := range sylType {
				var distr *distr.StrDistr
				switch c {
				case 'v':
					distr = distrVowels
				case 'c':
					distr = distrConsonants
				case 'l':
					distr = distrLetters
				default:
					log.Fatalf("main: unrecognized letter type %q\n", c)
				}
				buf.WriteString(distr.Pick())
			}
		}
		fmt.Println(strings.Title(buf.String()))
	}
}

// Msv - short for "must've"
func Msv(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	uniformVowels     = distr.NewStrDistr()
	uniformConsonants = distr.NewStrDistr()
	englishVowels     = distr.NewStrDistr()
	englishConsonants = distr.NewStrDistr()
)

func init() {
	uniformVowels.Add("a", 1)
	uniformConsonants.Add("b", 1)
	uniformConsonants.Add("c", 1)
	uniformConsonants.Add("d", 1)
	uniformVowels.Add("e", 1)
	uniformConsonants.Add("f", 1)
	uniformConsonants.Add("g", 1)
	uniformConsonants.Add("h", 1)
	uniformVowels.Add("i", 1)
	uniformConsonants.Add("j", 1)
	uniformConsonants.Add("k", 1)
	uniformConsonants.Add("l", 1)
	uniformConsonants.Add("m", 1)
	uniformConsonants.Add("n", 1)
	uniformVowels.Add("o", 1)
	uniformConsonants.Add("p", 1)
	uniformConsonants.Add("q", 1)
	uniformConsonants.Add("r", 1)
	uniformConsonants.Add("s", 1)
	uniformConsonants.Add("t", 1)
	uniformVowels.Add("u", 1)
	uniformConsonants.Add("v", 1)
	uniformConsonants.Add("w", 1)
	uniformConsonants.Add("x", 1)
	uniformVowels.Add("y", 1)
	uniformConsonants.Add("z", 1)

	englishVowels.Add("a", 8167)
	englishConsonants.Add("b", 1492)
	englishConsonants.Add("c", 2782)
	englishConsonants.Add("d", 4253)
	englishVowels.Add("e", 12702)
	englishConsonants.Add("f", 2228)
	englishConsonants.Add("g", 2015)
	englishConsonants.Add("h", 6094)
	englishVowels.Add("i", 6966)
	englishConsonants.Add("j", 153)
	englishConsonants.Add("k", 772)
	englishConsonants.Add("l", 4025)
	englishConsonants.Add("m", 2406)
	englishConsonants.Add("n", 6749)
	englishVowels.Add("o", 7507)
	englishConsonants.Add("p", 1929)
	englishConsonants.Add("q", 95)
	englishConsonants.Add("r", 5987)
	englishConsonants.Add("s", 6327)
	englishConsonants.Add("t", 9056)
	englishVowels.Add("u", 2758)
	englishConsonants.Add("v", 978)
	englishConsonants.Add("w", 2360)
	englishConsonants.Add("x", 150)
	englishVowels.Add("y", 1974)
	englishConsonants.Add("z", 74)
}
