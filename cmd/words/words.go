package words

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Metadata struct {
	Name       string
	Size       int
	PackagedAt string //Use some kind of date type here?
	Version    int
}

type WordSource struct {
	Metadata Metadata
	Words    []string
}

//go:embed embedables/words/common-english.json
var commonEnglish string

//go:embed embedables/sentences/frankenstein.json
var frankensteinSentences string

//go:embed embedables/words/korean-common.json
var koreanCommon string

//go:embed embedables/words/korean-tech.json
var koreanTech string

//go:embed embedables/sentences/korean-sentences.json
var koreanSentences string

func init() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
}

type WordsGenerator struct {
	Count     int
	pools     map[string]string
	poolsJson map[string]WordSource
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}



func addEmbededSource(sources map[string]WordSource) map[string]WordSource {
	var wordSource WordSource
	err := json.Unmarshal([]byte(commonEnglish), &wordSource)
	check(err)
	var sentenceSource WordSource
	err = json.Unmarshal([]byte(frankensteinSentences), &sentenceSource)
	check(err)

	sources["Common words"] = wordSource
	sources["Frankenstein sentences"] = sentenceSource

	return sources
}

func addEmbeddedKoreanSources(sources map[string]WordSource) map[string]WordSource {
	var koreanCommonSource WordSource
	err := json.Unmarshal([]byte(koreanCommon), &koreanCommonSource)
	check(err)
	
	var koreanTechSource WordSource
	err = json.Unmarshal([]byte(koreanTech), &koreanTechSource)
	check(err)
	
	var koreanSentenceSource WordSource
	err = json.Unmarshal([]byte(koreanSentences), &koreanSentenceSource)
	check(err)

	sources["Korean common words"] = koreanCommonSource
	sources["Korean tech terms"] = koreanTechSource
	sources["Korean sentences"] = koreanSentenceSource

	return sources
}

func unmarshalSources(paths []string) map[string]WordSource {
	acc := make(map[string]WordSource, len(paths))
	for _, sourceFilePath := range paths {
		var wordSource WordSource
		if strings.HasSuffix(sourceFilePath, ".json") {
			wordSource = readJsonSource(sourceFilePath)
		} else {
			wordSource = readNewLineSource(sourceFilePath)
		}

		acc[sourceFilePath] = wordSource
	}

	return acc
}

func readJsonSource(sourceFilePath string) WordSource {
	var wordSource WordSource

	fh, err := os.Open(sourceFilePath)
	defer fh.Close()
	check(err)

	decoder := json.NewDecoder(fh)
	err = decoder.Decode(&wordSource)
	check(err)

	return wordSource
}

func readNewLineSource(sourceFilePath string) WordSource {
	fh, err := os.Open(sourceFilePath)
	defer fh.Close()
	check(err)

	var lines []string
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	metadata := Metadata{
		Name:       fh.Name(),
		Size:       len(lines),
		PackagedAt: "1970-01-01T00:00:00Z",
		Version:    1,
	}

	return WordSource{
		Metadata: metadata,
		Words:    lines,
	}
}

func NewGenerator(paths []string) (g WordsGenerator) {
	g.Count = 300
	g.poolsJson = unmarshalSources(paths)
	g.poolsJson = addEmbededSource(g.poolsJson)

	return g
}

func NewKoreanGenerator(paths []string) (g WordsGenerator) {
	g.Count = 300
	g.poolsJson = unmarshalSources(paths)
	g.poolsJson = addEmbeddedKoreanSources(g.poolsJson)

	return g
}

func (this WordsGenerator) Generate(listName string) []rune {
	pool := this.poolsJson[listName].Words

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })

	takeAmount := min(this.Count, len(pool))
	words := pool[0:takeAmount]

	return []rune(strings.Join(words, " "))
}

// GenerateKoreanWords generates Korean words with proper shuffling and selection
func (this WordsGenerator) GenerateKoreanWords(listName string) []rune {
	pool := this.poolsJson[listName].Words

	// Korean-specific shuffling
	shuffledPool := this.shuffleKoreanWords(pool)
	
	takeAmount := min(this.Count, len(shuffledPool))
	words := shuffledPool[0:takeAmount]

	return []rune(strings.Join(words, " "))
}

// GenerateKoreanSentences generates Korean sentences with proper spacing and punctuation
func (this WordsGenerator) GenerateKoreanSentences(listName string) []rune {
	pool := this.poolsJson[listName].Words

	// Korean sentence generation with proper spacing
	shuffledPool := this.shuffleKoreanWords(pool)
	
	takeAmount := min(this.Count, len(shuffledPool))
	sentences := shuffledPool[0:takeAmount]

	// Join sentences with proper Korean spacing (sentences already have punctuation)
	return []rune(strings.Join(sentences, " "))
}

// shuffleKoreanWords performs Korean-specific word shuffling
func (this WordsGenerator) shuffleKoreanWords(words []string) []string {
	// Create a copy to avoid modifying the original slice
	shuffled := make([]string, len(words))
	copy(shuffled, words)
	
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(shuffled), func(i, j int) { 
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i] 
	})
	
	return shuffled
}

// GenerateKoreanWithCount generates Korean content respecting word/sentence count limits
func (this WordsGenerator) GenerateKoreanWithCount(listName string, count int) []rune {
	pool := this.poolsJson[listName].Words
	
	shuffledPool := this.shuffleKoreanWords(pool)
	
	// Respect the specified count limit
	takeAmount := min(count, len(shuffledPool))
	words := shuffledPool[0:takeAmount]
	
	return []rune(strings.Join(words, " "))
}
