package statesmanager

type GameState int

const (
	LoginState GameState = iota
	HCMenuState
	CreateUnitState
	InfoUnitState
	CreateUserState
	InfoUserState
	InboxState
	CreateDeviceState
	CreateTaskState
	GoBackState
)

type StateManager struct {
	Flow         chan GameState
	SceneStack   []GameState
	CurrentState GameState
}

func (s *StateManager) Add(state GameState) {
	s.Flow <- state
}
