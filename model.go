package angus

type Msg interface{}
type Html interface{}

type Model interface {
	Update(msg Msg) Model
	View() Html
}
