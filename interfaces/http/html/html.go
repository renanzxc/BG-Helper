package main

import (
	"github.com/labstack/echo/v4"
	"github.com/renanzxc/BG-Helper/bgg"
	bgghttp "github.com/renanzxc/BG-Helper/bgg/http"
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
		tmpl       *template.Template
		useCache   = true
		reqParams  struct {
			Username      string `param:"username" validate:"required"`
			AnotherPlayer string `param:"another_player"`
			ReloadCache   *bool  `query:"reload_cache"`
			NumPlayers    string `query:"num_players"`
		}
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
	sort.Slice(boardgames, func(i, j int) bool {
		return boardgames[i].Name < boardgames[j].Name
	})

	// TODO: read file only one time
	file, err := content.ReadFile("templates/playnext.html")
	if err != nil {
		return eris.Wrap(err, "")
	}
	if tmpl, err = template.New("playnext").Parse(string(file)); err != nil {
		return eris.Wrap(err, "")
	}

	return eris.Wrap(tmpl.Execute(writer, struct {
		Boardgames       []bgg.OwnedBoardgame
		AnotherPlayer    string
		NumPlayersFilter string
	}{
		Boardgames:       boardgames,
		AnotherPlayer:    reqParams.AnotherPlayer,
		NumPlayersFilter: reqParams.NumPlayers,
	}), "")
}
