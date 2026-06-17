package commandargs

type Args []string

func (a Args) IsEmpty() bool {
	return len(a) == 0
}

func (a Args) Len() int {
	return len(a)
}

func (a Args) First() string {
	return a[0]
}

func (a Args) Rest() Args {
	return a[1:]
}

func (a Args) At(index int) string {
	return a[index]
}
