package main

import (
	"github.com/labstack/echo/v4"
	"github.com/renanzxc/BG-Helper/bgg"
	bgghttp "github.com/renanzxc/BG-Helper/bgg/http"
	"github.com/rotisserie/eris"
	"io"
	"path"

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
	)

	if bgghttpImp, err = bgghttp.NewBGGClient(h.cache); err != nil {
		return
	}

	bggImp = bgg.NewBGG(bgghttpImp)

	boardgames, err := bggImp.GetBoardgamesToPlayNextWithAnotherUser(c.Request().Context(), c.Param("username"), c.Param("another_player"), true)
	if err != nil {
		return
	}
	sort.Slice(boardgames, func(i, j int) bool {
		return boardgames[i].Name < boardgames[j].Name
	})

	if tmpl, err = template.ParseFiles(path.Join(h.basePathTemplates, "playnext.html")); err != nil {
		return eris.Wrap(err, "")
	}

	return eris.Wrap(tmpl.Execute(writer, struct {
		Boardgames []bgg.OwnedBoardgame
	}{
		Boardgames: boardgames,
	}), "")
}