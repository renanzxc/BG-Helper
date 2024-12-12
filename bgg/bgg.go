package bgg

import (
	"context"
	"github.com/renanzxc/BG-Helper/bgg/http"
	"math"
)

type BGG interface {
	GetBoardgamesToPlayNextWithAnotherUser(ctx context.Context, username, anotherUsername string, useCache bool) (boardgames []OwnedBoardgame, err error)
}

type BGGimp struct {
	bggHTTP http.HTTP
}

func NewBGG(http http.HTTP) *BGGimp {
	return &BGGimp{bggHTTP: http}
}

func (b *BGGimp) GetBoardgamesToPlayNextWithAnotherUser(ctx context.Context, username, anotherUsername string, useCache bool) (boardgames []OwnedBoardgame, err error) {
	col1, err := b.GetUserCollection(ctx, username, useCache)
	if err != nil {
		return nil, err
	}
	col2, err := b.GetUserCollection(ctx, anotherUsername, useCache)
	if err != nil {
		return nil, err
	}
	j1 := col1.WantToPlayInOwnedCollection(col2)
	j2 := col2.WantToPlayInOwnedCollection(col1)
	boardgames = append(j1, j2...)

	for ii := range boardgames {
		// Always use boardgame cache
		if err = b.MoreInfo(ctx, &boardgames[ii].Boardgame, true); err != nil {
			return
		}
	}

	return
}

func (b *BGGimp) GetUserCollection(ctx context.Context, username string, useCache bool) (collection *UserCollection, err error) {
	var colXML *http.CollectionXML

	if colXML, err = b.bggHTTP.GetUserCollection(ctx, username, useCache); err != nil {
		return
	}
	collection = &UserCollection{
		Username:   username,
		Owned:      make(map[string]Boardgame),
		Preordered: make(map[string]Boardgame),
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

		if item.Status.Preordered == 1 {
			collection.Preordered[item.ObjectID] = Boardgame{
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
	boardgame.MaxPlayTime = thing.Items[0].MaxPlayTime.Value
	boardgame.MinPlayTime = thing.Items[0].MinPlayTime.Value
	boardgame.AverageWeight = math.Round(thing.Items[0].Statistics[0].Ratings[0].AverageWeight.Value*100) / 100

	return
}
