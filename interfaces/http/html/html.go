package main

import (
	"github.com/labstack/echo/v4"
	"github.com/renanzxc/BG-Helper/bgg"
	bgghttp "github.com/renanzxc/BG-Helper/bgg/http"
	"github.com/renanzxc/BG-Helper/watcher"
	"github.com/rotisserie/eris"
	"io"
	"sort"
	"text/template"
)

func playnext(c echo.Context, h *HTTPHTML) error {
	err := playNextHTML(c, h, c.Response().Writer)
	if err != nil {
		return err
	}

	return nil
}

func playNextHTML(c echo.Context, h *HTTPHTML, writer io.Writer) (err error) {
	var (
		bgghttpImp bgghttp.HTTP
		bggImp     bgg.BGG
		tmplMain   *template.Template
		useCache   = true
		reqParams  struct {
			Username      string `param:"username" validate:"required"`
			AnotherPlayer string `param:"another_player"`
			ReloadCache   *bool  `query:"reload_cache"`
			NumPlayers    string `query:"num_players"`
		}
		inProgress bool
	)

	if err = c.Bind(&reqParams); err != nil {
		return eris.Wrap(err, "")
	}

	if err = h.validate.Struct(&reqParams); err != nil {
		return eris.Wrap(err, "")
	}

	if reqParams.ReloadCache != nil {
		useCache = !*reqParams.ReloadCache
	}

	if bgghttpImp, err = bgghttp.NewBGGClient(h.cache); err != nil {
		return
	}

	bggImp = bgg.NewBGG(bgghttpImp)

	boardgames, err := bggImp.GetBoardgamesToPlayNextWithAnotherUser(c.Request().Context(), reqParams.Username, reqParams.AnotherPlayer, useCache)
	if err != nil {
		return
	}

	inProgress = boardgames.State == watcher.InProgress

	if boardgames.State == watcher.Done {
		sort.Slice(boardgames.Boardgames, func(i, j int) bool {
			return boardgames.Boardgames[i].Name < boardgames.Boardgames[j].Name
		})
	}

	// TODO: read playNextHTMLFile only one time
	playNextHTMLFile, err := content.ReadFile("templates/playnext.html")
	if err != nil {
		return eris.Wrap(err, "")
	}

	if tmplMain, err = template.New("playnext").Parse(string(playNextHTMLFile)); err != nil {
		return eris.Wrap(err, "")
	}

	var loadingHTMLFile = []byte("")

	if inProgress {
		loadingHTMLFile, err = content.ReadFile("templates/loading.html")
		if err != nil {
			return eris.Wrap(err, "")
		}
	}

	if tmplMain, err = tmplMain.New("loading").Parse(string(loadingHTMLFile)); err != nil {
		return eris.Wrap(err, "")
	}

	return eris.Wrap(tmplMain.ExecuteTemplate(writer, "playnext", struct {
		Boardgames       []bgg.OwnedBoardgame
		AnotherPlayer    string
		NumPlayersFilter string
		InProgress       bool
	}{
		Boardgames:       boardgames.Boardgames,
		AnotherPlayer:    reqParams.AnotherPlayer,
		NumPlayersFilter: reqParams.NumPlayers,
		InProgress:       inProgress,
	}), "")
}
