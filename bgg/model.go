package bgg

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
	OwnedByUsername string
}
