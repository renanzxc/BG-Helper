package bgg

type UserCollection struct {
	Username   string
	Owned      map[string]Boardgame
	Preordered map[string]Boardgame
	WantToPlay map[string]Boardgame
}

func (col *UserCollection) WantToPlayInMyCollection() (boardgames []Boardgame) {
	for boardgameID := range col.WantToPlay {
		boardgames = append(boardgames, col.WantToPlay[boardgameID])
	}

	return
}

func (col *UserCollection) WantToPlayInOwnedCollection(colIn *UserCollection) (boardgames []OwnedBoardgame) {
	var addedBoardgames = make(map[string]bool)

	for boardgameID, wantToPlayData := range col.WantToPlay {
		if wantToPlayData.NumPlays != 0 {
			continue
		}
		boardgameColIn, ok := colIn.Owned[boardgameID]
		if !ok {
			boardgameColIn, ok = colIn.Preordered[boardgameID]
			if !ok {
				continue
			}
		}

		if addedBoardgames[boardgameID] {
			continue
		}

		addedBoardgames[boardgameID] = true
		boardgames = append(boardgames, OwnedBoardgame{
			Boardgame:             boardgameColIn,
			OwnedByUsername:       colIn.Username,
			OwndedNotPlayThisGame: boardgameColIn.NumPlays == 0,
		})
	}

	return
}
