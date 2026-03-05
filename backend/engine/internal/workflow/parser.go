package workflow

type Definition struct {
  Name string
}

func ParseYAML(_ []byte) (Definition, error) {
  return Definition{Name: "bootstrap"}, nil
}
