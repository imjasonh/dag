package pkg

<<<<<<< HEAD
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
=======
import "chainguard.dev/melange/pkg/build"

type URI struct {
	URI, ExpectedSHA256, ExpectedSHA512 string
}

func URIs(cfg build.Configuration) []URI {
	var uris []URI
	for _, s := range cfg.Pipeline {
		if s.Uses == "fetch" {
			uris = append(uris, URI{
				URI:            s.With["uri"],
				ExpectedSHA256: s.With["expected-sha256"],
				ExpectedSHA512: s.With["expected-sha512"],
			})
		}
	}
	return uris
>>>>>>> baa37ce (stash)
}
