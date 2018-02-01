package handlers

import (
	"github.com/labstack/echo"
	"github.com/oleggator/esports-backend/db"
	"github.com/oleggator/esports-backend/models"
	"github.com/jackc/pgx"
	"net/http"
	"io/ioutil"
	"github.com/Machiel/slugify"
	"fmt"
)

func GetGames(ctx echo.Context) (err error) {
	since := ctx.Get("since").(int64)
	limit := ctx.Get("limit").(int32)
	if limit == 0 || limit > 100 {
		limit = 100
	}
	desc := ctx.Get("desc").(bool)

	var rows *pgx.Rows
	if since == 0 && !desc {
		rows, err = db.Read(0).Query("selectAllGamesAsc", limit)
	} else if since == 0 && desc {
		rows, err = db.Read(0).Query("selectAllGamesDesc", limit)
	} else if since != 0 && !desc {
		rows, err = db.Read(0).Query("selectAllGamesAscSince", since, limit)
	} else {
		rows, err = db.Read(0).Query("selectAllGamesDescSince", since, limit)
	}
	defer rows.Close()
	if err != nil {
		ctx.Logger().Error(err)
	}

	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	ctx.Response().WriteHeader(http.StatusOK)

	first := true
	ctx.Response().Write([]byte("["))
	for rows.Next() {
		if first {
			first = false
		} else {
			ctx.Response().Write([]byte(","))
		}

		game := models.Game{}
		err = rows.Scan(&game.ID, &game.Title, &game.Slug)
		if err != nil {
			ctx.Logger().Error(err)
		}

		blob, err := game.MarshalBinary()
		if err != nil {
			ctx.Logger().Error(err)
		}
		ctx.Response().Write(blob)
	}

	ctx.Response().Write([]byte("]"))
	ctx.Response().Flush()

	return nil
}

func GetGame(ctx echo.Context) (err error) {
	slug := ctx.Get("slug").(string)
	shardIndex := ctx.Get("shardIndex").(int)

	game := models.Game{}
	err = db.Read(shardIndex).QueryRow("selectGame", slug).Scan(&game.ID, &game.Title, &game.Slug)
	if err != nil {
		ctx.Logger().Error(err)
		apiErr := models.Error{Message: "Not found"}
		blob, err := apiErr.MarshalBinary()
		if err != nil {
			ctx.Logger().Error(err)
		}

		ctx.JSONBlob(http.StatusNotFound, blob)

		return err
	}

	blob, err := game.MarshalBinary()
	if err != nil {
		ctx.Logger().Error(err)
	}

	return ctx.JSONBlob(http.StatusOK, blob)
}

func AddGame(ctx echo.Context) (err error) {
	bodyBlob, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		ctx.Logger().Error(err)
	}

	game := models.Game{}
	err = game.UnmarshalBinary(bodyBlob)
	if err != nil {
		ctx.Logger().Error(err)
	}

	const shardIndex = 0
	tx, err := db.Write(shardIndex).Begin()

	game.Slug = fmt.Sprintf("%d-%s", shardIndex, slugify.Slugify(*game.Title))
	err = tx.QueryRow("insertGame", game.Title, game.Slug).Scan(&game.ID)
	if err != nil {
		ctx.Logger().Error(err)

		apiErr := models.Error{Message: err.Error()}
		blob, err := apiErr.MarshalBinary()
		if err != nil {
			ctx.Logger().Error(err)
		}

		err = tx.Rollback()
		if err != nil {
			ctx.Logger().Error(err)
		}

		return ctx.JSONBlob(http.StatusConflict, blob)
	}

	err = tx.Commit()
	if err != nil {
		ctx.Logger().Error(err)
	}

	blob, err := game.MarshalBinary()
	if err != nil {
		ctx.Logger().Error(err)
	}

	return ctx.JSONBlob(http.StatusCreated, blob)
}

func DeleteGame(ctx echo.Context) (err error) {
	slug := ctx.Get("slug").(string)
	shardIndex := ctx.Get("shardIndex").(int)

	tx, err := db.Write(shardIndex).Begin()
	_, err = tx.Exec("deleteGame", slug)
	if err != nil {
		ctx.Logger().Error(err)

		apiError := models.Error{Message: "Game not found"}
		blob, err := apiError.MarshalBinary()
		if err != nil {
			ctx.Logger().Error(err)
		}

		err = tx.Rollback()
		if err != nil {
			ctx.Logger().Error(err)
		}

		return ctx.JSONBlob(http.StatusNotFound, blob)
	}
	tx.Commit()

	return ctx.NoContent(http.StatusOK)
}
