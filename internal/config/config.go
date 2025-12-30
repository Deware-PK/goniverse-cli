package config

type Config struct {
	EpubPath    string
	OutputDir   string
	AllowedTags map[string]bool
	Workers     int
}

func NewDefaultConfig() *Config {
	return &Config{
		Workers: 10,
		AllowedTags: map[string]bool{
			"p":          true,
			"h1":         true,
			"h2":         true,
			"h3":         true,
			"h4":         true,
			"h5":         true,
			"h6":         true,
			"b":          true,
			"i":          true,
			"strong":     true,
			"em":         true,
			"u":          true,
			"br":         true,
			"div":        true,
			"span":       true,
			"li":         true,
			"ul":         true,
			"ol":         true,
			"blockquote": true,
			"img":        true,
		},
	}
}
