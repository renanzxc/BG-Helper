package bgg

type UserCollection struct {
	Username   string
	Owned      map[string]Boardgame
	WantToPlay map[string]Boardgame
}

func (col *UserCollection) WantToPlayInMyCollection() (boardgames []Boardgame) {
	for boardgameID := range col.WantToPlay {
		boardgames = append(boardgames, col.WantToPlay[boardgameID])
	}

	return
}

func (col *UserCollection) WantToPlayInOwnedCollection(colIn *UserCollection) (boardgames []OwnedBoardgame) {
	for boardgameID, wantToPlayData := range col.WantToPlay {
		if wantToPlayData.NumPlays != 0 {
			continue
		}
		boardgameColIn, ok := colIn.Owned[boardgameID]
		if ok {
			boardgames = append(boardgames, OwnedBoardgame{
				Boardgame:             boardgameColIn,
				OwnedByUsername:       colIn.Username,
				OwndedNotPlayThisGame: boardgameColIn.NumPlays == 0,
			})
		}
	}

	return
}
