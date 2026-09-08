package engine

import (
	"math/rand/v2"
	"regexp"
	"strings"
	"unicode"
)

type Quote struct {
	Text     string   `json:"text"`
	Category string   `json:"category"`
	Triggers []string `json:"triggers"`
	Weight   int      `json:"weight"`
}

type Result struct {
	Reply      bool
	Text       string
	Score      float64
	Summoned   bool
	Category   string
	MatchedKey string
	Chance     float64
}

var spaces = regexp.MustCompile(`\s+`)

func DefaultSummonWords() []string {
	return []string{"野獣先輩", "やじゅ", "やじゅせん", "yaju", "yajuu", "114514", "810", "淫夢"}
}

func Normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		// Full-width ASCII -> ASCII, ideographic space -> regular space.
		if r >= 0xFF01 && r <= 0xFF5E {
			r -= 0xFEE0
		} else if r == 0x3000 {
			r = ' '
		}
		switch {
		case unicode.IsLetter(r), unicode.IsNumber(r):
			b.WriteRune(r)
		case unicode.IsSpace(r):
			b.WriteRune(' ')
		}
	}
	return spaces.ReplaceAllString(strings.TrimSpace(b.String()), " ")
}

func isASCIIDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func containsDigitToken(message, token string) bool {
	norm := Normalize(message)
	var run strings.Builder
	flush := func() bool {
		if run.String() == token {
			return true
		}
		run.Reset()
		return false
	}
	for _, r := range norm {
		if r >= '0' && r <= '9' {
			run.WriteRune(r)
			continue
		}
		if flush() {
			return true
		}
	}
	return flush()
}

func containsNormalized(haystack, needle string) bool {
	h := Normalize(haystack)
	n := Normalize(needle)
	if n == "" {
		return false
	}
	if isASCIIDigits(n) {
		return containsDigitToken(haystack, n)
	}
	return strings.Contains(h, n)
}

func modeMultiplier(mode string) float64 {
	switch strings.ToLower(mode) {
	case "quiet":
		return 0.45
	case "aggressive":
		return 1.6
	case "chaos":
		return 2.5
	default:
		return 1.0
	}
}

func Evaluate(message string, summonWords []string, baseChance float64, mode string, quotes []Quote) Result {
	norm := Normalize(message)
	if norm == "" || len(quotes) == 0 {
		return Result{}
	}

	summoned := false
	matchedKey := ""
	for _, w := range summonWords {
		if containsNormalized(message, w) {
			summoned = true
			matchedKey = w
			break
		}
	}

	scores := map[string]float64{}
	for _, q := range quotes {
		for _, t := range q.Triggers {
			if containsNormalized(message, t) {
				lengthBonus := float64(len([]rune(Normalize(t)))) / 10.0
				scores[q.Category] += 1.0 + lengthBonus
			}
		}
	}

	bestCategory := "fallback"
	bestScore := 0.0
	for category, score := range scores {
		if score > bestScore {
			bestCategory = category
			bestScore = score
		}
	}

	candidates := make([]Quote, 0, 8)
	for _, q := range quotes {
		if q.Category == bestCategory || (bestScore == 0 && q.Category == "fallback") {
			candidates = append(candidates, q)
		}
	}
	if len(candidates) == 0 {
		candidates = quotes
	}

	chance := baseChance
	if summoned {
		chance = 100
		bestScore += 10
	} else {
		chance = (chance + bestScore*7.0) * modeMultiplier(mode)
	}
	if chance > 100 {
		chance = 100
	}
	if chance < 0 {
		chance = 0
	}

	result := Result{
		Score:      bestScore,
		Summoned:   summoned,
		Category:   bestCategory,
		MatchedKey: matchedKey,
		Chance:     chance,
	}
	if rand.Float64()*100 >= chance {
		return result
	}

	chosen := weightedPick(candidates)
	result.Reply = true
	result.Text = chosen.Text
	return result
}

func weightedPick(quotes []Quote) Quote {
	total := 0
	for _, q := range quotes {
		w := q.Weight
		if w <= 0 {
			w = 1
		}
		total += w
	}
	pick := rand.IntN(max(total, 1))
	for _, q := range quotes {
		w := q.Weight
		if w <= 0 {
			w = 1
		}
		if pick < w {
			return q
		}
		pick -= w
	}
	return quotes[0]
}

func RandomQuote(quotes []Quote) string {
	if len(quotes) == 0 {
		return ""
	}
	return weightedPick(quotes).Text
}
