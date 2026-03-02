package templates

type Generator interface {
	Name() string
	Generate(projectPath, modName string) error
	PrintNextSteps()
}

func Get(name string) (Generator, bool) {
	return getGenerator(name)
}

func List() []string {
	return listGenerators()
}
