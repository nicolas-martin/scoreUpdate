package Interfaces

// GameEvents that all games should have
type GameEvents interface {
	StartGame() bool
	Goal() bool
	EndPeriod() bool
	EndGame() bool
	Loop()
}

//Game struct that all games should have
type Game struct {
	ID  int
	URL string
}
