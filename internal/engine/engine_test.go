package engine

import "testing"

var testQuotes = []Quote{
    {Text:"PRAISE", Category:"praise", Triggers:[]string{"すごい","成功"}, Weight:1},
    {Text:"TIRED", Category:"tired", Triggers:[]string{"疲れ","眠い"}, Weight:1},
    {Text:"FALLBACK", Category:"fallback", Weight:1},
}

func TestNormalize(t *testing.T) {
    got := Normalize("  Yaju！！ 先輩 ")
    if got != "yaju 先輩" { t.Fatalf("Normalize=%q", got) }
}

func TestSummonAlwaysReplies(t *testing.T) {
    r := Evaluate("野獣先輩いる？", DefaultSummonWords(), 0, "quiet", testQuotes)
    if !r.Reply || !r.Summoned { t.Fatalf("expected summoned reply: %+v", r) }
}

func TestContextCategory(t *testing.T) {
    r := Evaluate("今日めっちゃ疲れた", nil, 100, "normal", testQuotes)
    if !r.Reply || r.Category != "tired" || r.Text != "TIRED" { t.Fatalf("unexpected result: %+v", r) }
}

func TestZeroChanceNoContext(t *testing.T) {
    r := Evaluate("こんにちは", nil, 0, "normal", testQuotes)
    if r.Reply { t.Fatalf("unexpected reply: %+v", r) }
}
