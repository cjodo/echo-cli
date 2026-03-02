package templates

type generatorRegistry map[string]Generator

var registry generatorRegistry = make(map[string]Generator)

// Register adds a template generator to the registry.
func Register(g Generator) {
	registry[g.Name()] = g
}

func getGenerator(name string) (Generator, bool) {
	g, ok := registry[name]
	return g, ok
}

func listGenerators() []string {
	keys := make([]string, 0, len(registry))
	for k := range registry {
		keys = append(keys, k)
	}
	return keys
}
