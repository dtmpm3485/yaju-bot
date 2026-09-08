package store

import (
	"path/filepath"
	"testing"
)

func TestStoreRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Update("guild", func(c *GuildConfig) {
		c.Chance = 12.5
		c.Mode = "aggressive"
	}); err != nil {
		t.Fatal(err)
	}
	s.UpdateVolatile("guild", func(c *GuildConfig) { c.Seen++ })
	if err := s.Flush(); err != nil {
		t.Fatal(err)
	}

	s2, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	got := s2.Get("guild")
	if got.Chance != 12.5 || got.Mode != "aggressive" || got.Seen != 1 {
		t.Fatalf("unexpected config: %+v", got)
	}
}
