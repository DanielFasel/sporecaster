package spore

type File struct {
	Name string `yaml:"name" json:"name"`
	Role string `yaml:"role" json:"role"`
}

type Core struct {
	Name        string   `yaml:"name"        json:"name"`
	Kind        string   `yaml:"kind"        json:"kind"`
	Description string   `yaml:"description" json:"description"`
	Files       []File   `yaml:"files"       json:"files"`
	Imports     []string `yaml:"imports"     json:"imports"`
}

type Field struct {
	Name string `yaml:"name" json:"name"`
	Type string `yaml:"type" json:"type"`
}

type Method struct {
	Name      string `yaml:"name"      json:"name"`
	Signature string `yaml:"signature" json:"signature"`
}

type Export struct {
	Name      string   `yaml:"name"      json:"name"`
	Kind      string   `yaml:"kind"      json:"kind"`
	Signature string   `yaml:"signature" json:"signature"`
	Fields    []Field  `yaml:"fields"    json:"fields"`
	Methods   []Method `yaml:"methods"   json:"methods"`
}

type Package struct {
	Name        string   `yaml:"name"        json:"name"`
	Parent      string   `yaml:"parent"      json:"parent"`
	Kind        string   `yaml:"kind"        json:"kind"`
	Description string   `yaml:"description" json:"description"`
	Files       []File   `yaml:"files"       json:"files"`
	Imports     []string `yaml:"imports"     json:"imports"`
	Exports     []Export `yaml:"exports"     json:"exports"`
}

type ErrorHandling struct {
	TerminatesAt         string `yaml:"terminates_at"          json:"terminates_at"`
	SentinelsAreExported bool   `yaml:"sentinels_are_exported" json:"sentinels_are_exported"`
}

type Command struct {
	Name        string `yaml:"name"        json:"name"`
	Usage       string `yaml:"usage"       json:"usage"`
	Description string `yaml:"description" json:"description"`
}

type Channel struct {
	Name        string    `yaml:"name"        json:"name"`
	Type        string    `yaml:"type"        json:"type"`
	Description string    `yaml:"description" json:"description"`
	Commands    []Command `yaml:"commands"    json:"commands"`
}

type Spore struct {
	App           string        `yaml:"app"            json:"app"`
	Description   string        `yaml:"description"    json:"description"`
	Language      string        `yaml:"language"       json:"language"`
	ErrorHandling ErrorHandling `yaml:"error_handling" json:"error_handling"`
	Core          Core          `yaml:"core"           json:"core"`
	Packages      []Package     `yaml:"packages"       json:"packages"`
	Channels      []Channel     `yaml:"channels"       json:"channels"`
}
