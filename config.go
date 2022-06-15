package xmap

const (
	// DefaultInitialCap is the default for Config.InitialCap
	DefaultInitialCap = 8

	// DefaultGrowLoadLimit is the default for Config.GrowLoadLimit
	DefaultGrowLoadLimit = 0.8

	// DefaultShrinkLoadLimit is the default for Config.ShrinkLoadLimit
	DefaultShrinkLoadLimit = 0.25

	// DefaultChangeFactor is the default for Config.ChangeFactor
	DefaultChangeFactor = 2.0
)

// Config holds properties that may be changed.
// Fields holding the zero values of their types will be interpreted as defaults.
type Config struct {
	// InitialCap is the initial capacity to use when the first entry is added.
	// Checked when size increases and current capacity is 0.
	// Capacity will not be decreased below this automatically.
	InitialCap int

	// GrowLoadLimit is the load factor over which entries are rehashed.
	// Checked when capacity is adjusted.
	GrowLoadLimit float64

	// ShrinkLoadLimit is the load factor under which entries are rehashed.
	// Checked when capacity is adjusted.
	ShrinkLoadLimit float64

	// ChangeFactor is the factor to calculate new capacity on increment or on shrink.
	ChangeFactor float64
}

// setDefaults sets default values for fields that are the zero values of their types.
func (cfg *Config) setDefaults() *Config {
	if cfg.InitialCap == 0 {
		cfg.InitialCap = DefaultInitialCap
	}
	if cfg.GrowLoadLimit == 0 {
		cfg.GrowLoadLimit = DefaultGrowLoadLimit
	}
	if cfg.ShrinkLoadLimit == 0 {
		cfg.ShrinkLoadLimit = DefaultShrinkLoadLimit
	}
	if cfg.ChangeFactor == 0 {
		cfg.ChangeFactor = DefaultChangeFactor
	}

	return cfg
}

var defaultConfig = (&Config{}).setDefaults()
