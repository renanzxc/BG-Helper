package bgg

type Boardgame struct {
	ObjectID            string
	Name                string
	NumPlays            int
	MinPlayers          int
	MaxPlayers          int
	BestNumPlayers      string
	BestNumPlayersVotes int
}

type OwnedBoardgame struct {
	Boardgame
	OwnedByUsername string
}
