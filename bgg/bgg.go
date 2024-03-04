package bgg

import (
	"context"
	"github.com/renanzxc/BG-Helper/bgg/http"
)

type BGG interface {
	MoreInfo(ctx context.Context, boardgame *Boardgame, useCache bool) (err error)
}

type BGGimp struct {
	bggHTTP http.HTTP
}

func NewBGG(http http.HTTP) *BGGimp {
	return &BGGimp{bggHTTP: http}
}

func (b *BGGimp) GetUserCollection(ctx context.Context, username string, useCache bool) (collection *UserCollection, err error) {
	var colXML *http.CollectionXML

	if colXML, err = b.bggHTTP.GetUserCollection(ctx, username, useCache); err != nil {
		return
	}
	collection = &UserCollection{
		Username:   username,
		Owned:      make(map[string]Boardgame),
		WantToPlay: make(map[string]Boardgame),
	}

	for _, item := range colXML.Items {
		if item.Status.Own == 1 {
			collection.Owned[item.ObjectID] = Boardgame{
				ObjectID: item.ObjectID,
				Name:     item.Name,
				NumPlays: item.NumPlays,
			}
		}

		if item.Status.WantToPlay == 1 {
			collection.WantToPlay[item.ObjectID] = Boardgame{
				ObjectID: item.ObjectID,
				Name:     item.Name,
				NumPlays: item.NumPlays,
			}
		}
	}

	return
}

func (b *BGGimp) MoreInfo(ctx context.Context, boardgame *Boardgame, useCache bool) (err error) {
	var (
		thing        *http.ThingsXML
		best         *http.ThingResultsXML
		numVotesBest = -1
	)

	if thing, err = b.bggHTTP.GetThing(ctx, boardgame.ObjectID, useCache); err != nil {
		return
	}

	for _, poll := range thing.Items[0].Poll {
		if poll.Name == "suggested_numplayers" {
			for ii := range poll.Results {
				if numVotes := poll.Results[ii].GetBestNumVotes(); numVotes > numVotesBest {
					best = &poll.Results[ii]
					numVotesBest = numVotes
				}
			}
		}
	}

	boardgame.MinPlayers = thing.Items[0].MinPlayers.Value
	boardgame.MaxPlayers = thing.Items[0].MaxPlayers.Value
	boardgame.BestNumPlayers = best.NumPlayers
	boardgame.BestNumPlayersVotes = numVotesBest

	return
}
