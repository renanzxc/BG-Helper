package bgg

import "github.com/renanzxc/BG-Helper/watcher"

type Boardgame struct {
	ObjectID            string
	Name                string
	NumPlays            int
	MinPlayers          int
	MaxPlayers          int
	MinPlayTime         int
	MaxPlayTime         int
	BestNumPlayers      string
	BestNumPlayersVotes int
	AverageWeight       float64
}

type OwnedBoardgame struct {
	Boardgame
	OwnedByUsername       string
	OwndedNotPlayThisGame bool
}

type OwnedBoardgames struct {
	watcher.ProcessReturn
	Boardgames []OwnedBoardgame
}
