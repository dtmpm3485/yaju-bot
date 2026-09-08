package engine

import (
    "math/rand/v2"
    "regexp"
    "strings"
    "unicode"
)

type Quote struct {
    Text string `json:"text"`
    Category string `json:"category"`
    Triggers []string `json:"triggers"`
    Weight int `json:"weight"`
}

type Result struct {
    Reply bool
    Text string
    Score float64
    Summoned bool
    Category string
    MatchedKey string
}

var spaces = regexp.MustCompile(`\s+`)

func DefaultSummonWords() []string {
    return []string{"野獣先輩", "やじゅ", "yaju", "yajuu", "114514", "810"}
}

func Normalize(s string) string {
    s = strings.ToLower(strings.TrimSpace(s))
    var b strings.Builder
    b.Grow(len(s))
    for _, r := range s {
        switch {
        case unicode.IsLetter(r), unicode.IsNumber(r):
            b.WriteRune(r)
        case unicode.IsSpace(r):
            b.WriteRune(' ')
        }
    }
    return spaces.ReplaceAllString(b.String(), " ")
}

func containsNormalized(haystack, needle string) bool {
    h := Normalize(haystack)
    n := Normalize(needle)
    return n != "" && strings.Contains(h, n)
}

func Evaluate(message string, summonWords []string, baseChance float64, mode string, quotes []Quote) Result {
    norm := Normalize(message)
    if norm == "" || len(quotes) == 0 { return Result{} }

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
    bestCategory := "fallback"
    bestScore := 0.0
    for _, q := range quotes {
        for _, t := range q.Triggers {
            if containsNormalized(message, t) {
                scores[q.Category] += 1.0 + float64(len([]rune(Normalize(t))))/12.0
            }
        }
        if scores[q.Category] > bestScore {
            bestScore = scores[q.Category]
            bestCategory = q.Category
        }
    }

    candidates := make([]Quote, 0, 8)
    for _, q := range quotes {
        if q.Category == bestCategory || (bestScore == 0 && q.Category == "fallback") {
            candidates = append(candidates, q)
        }
    }
    if len(candidates) == 0 { candidates = quotes }

    chance := baseChance
    if summoned {
        chance = 100
        bestScore += 10
    } else {
        chance += bestScore * 8
        switch strings.ToLower(mode) {
        case "quiet": chance *= 0.45
        case "chaos": chance *= 2.2
        }
    }
    if chance > 100 { chance = 100 }
    if chance < 0 { chance = 0 }
    if rand.Float64()*100 >= chance {
        return Result{Score:bestScore, Summoned:summoned, Category:bestCategory, MatchedKey:matchedKey}
    }

    total := 0
    for _, q := range candidates { if q.Weight > 0 { total += q.Weight } else { total++ } }
    pick := rand.IntN(max(total, 1))
    chosen := candidates[0]
    for _, q := range candidates {
        w := q.Weight; if w <= 0 { w = 1 }
        if pick < w { chosen = q; break }
        pick -= w
    }
    return Result{Reply:true, Text:chosen.Text, Score:bestScore, Summoned:summoned, Category:bestCategory, MatchedKey:matchedKey}
}

func RandomQuote(quotes []Quote) string {
    if len(quotes)==0 { return "" }
    return quotes[rand.IntN(len(quotes))].Text
}
