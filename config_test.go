package xmap

import (
	"testing"

	"github.com/icza/mighty"
)

var defConfig *Config = &Config{
	InitialCap:      DefaultInitialCap,
	GrowLoadLimit:   DefaultGrowLoadLimit,
	ShrinkLoadLimit: DefaultShrinkLoadLimit,
	ChangeFactor:    DefaultChangeFactor,
}

func TestConfig(t *testing.T) {
	deq := mighty.Deq(t)

	cfg := &Config{}

	cfg.setDefaults()
	deq(defConfig, cfg)

	refConfig := &Config{
		InitialCap:      1,
		GrowLoadLimit:   0.1,
		ShrinkLoadLimit: 0.2,
		ChangeFactor:    0.3,
	}

	*cfg = *refConfig
	cfg.setDefaults()
	deq(refConfig, cfg)
}
