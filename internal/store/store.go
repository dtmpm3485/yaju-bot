package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/dtmpm3485/yaju-bot/internal/engine"
)

type GuildConfig struct {
	Enabled         bool     `json:"enabled"`
	Chance          float64  `json:"chance"`
	CooldownSeconds int      `json:"cooldown_seconds"`
	Mode            string   `json:"mode"`
	Channels        []string `json:"channels"`
	SummonWords     []string `json:"summon_words"`
	Seen            uint64   `json:"seen"`
	Replies         uint64   `json:"replies"`
	Summons         uint64   `json:"summons"`
}

type Data struct {
	Guilds map[string]*GuildConfig `json:"guilds"`
}

type Store struct {
	mu    sync.RWMutex
	path  string
	data  Data
	dirty bool
}

func DefaultConfig() *GuildConfig {
	return &GuildConfig{
		Enabled:         true,
		Chance:          2.5,
		CooldownSeconds: 90,
		Mode:            "normal",
		SummonWords:     engine.DefaultSummonWords(),
	}
}

func Open(path string) (*Store, error) {
	s := &Store{path: path, data: Data{Guilds: map[string]*GuildConfig{}}}
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return s, nil
		}
		return nil, err
	}
	if len(b) > 0 {
		if err := json.Unmarshal(b, &s.data); err != nil {
			return nil, err
		}
	}
	if s.data.Guilds == nil {
		s.data.Guilds = map[string]*GuildConfig{}
	}
	return s, nil
}

func (s *Store) Get(guildID string) GuildConfig {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.data.Guilds[guildID]
	if !ok {
		c = DefaultConfig()
		s.data.Guilds[guildID] = c
		s.dirty = true
		_ = s.saveLocked()
	}
	cp := *c
	cp.Channels = append([]string(nil), c.Channels...)
	cp.SummonWords = append([]string(nil), c.SummonWords...)
	return cp
}

// Update changes configuration and persists immediately.
func (s *Store) Update(guildID string, fn func(*GuildConfig)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c := s.getOrCreateLocked(guildID)
	fn(c)
	s.dirty = true
	return s.saveLocked()
}

// UpdateVolatile is for high-frequency counters. It is flushed periodically.
func (s *Store) UpdateVolatile(guildID string, fn func(*GuildConfig)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c := s.getOrCreateLocked(guildID)
	fn(c)
	s.dirty = true
}

func (s *Store) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.dirty {
		return nil
	}
	return s.saveLocked()
}

func (s *Store) getOrCreateLocked(guildID string) *GuildConfig {
	c, ok := s.data.Guilds[guildID]
	if !ok {
		c = DefaultConfig()
		s.data.Guilds[guildID] = c
	}
	return c
}

func (s *Store) saveLocked() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	if err := os.Rename(tmp, s.path); err != nil {
		return err
	}
	s.dirty = false
	return nil
}
