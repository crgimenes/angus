package angus

type Msg interface{}
type Html interface{}

type Model interface {
	Update(msg Msg) Model
	View() Html
}

type Program struct {
	model Model
}

func NewProgram(init Model) Program {
	return Program{
		model: init,
	}
}