package pkg

// TODO: reuse melange's pkg/build.Configuration type
type Config struct {
	Package struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
		Epoch   string `yaml:"epoch"`
	}
	Environment struct {
		Contents struct {
			Packages []string
		}
	}
	Pipeline []struct {
		Uses string
		With map[string]string
	}
	Data []struct {
		Name  string
		Items map[string]string
	}
	Subpackages []Subpackage
}

type Subpackage struct {
	Name, Range string
}
