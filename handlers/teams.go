package handlers

import (
	"github.com/labstack/echo"
	"github.com/oleggator/esports-backend/db"
	"github.com/oleggator/esports-backend/models"
	"github.com/jackc/pgx"
	"net/http"
	"io/ioutil"
	"hash/adler32"
	"github.com/Machiel/slugify"
	"fmt"
)

func GetTeams(ctx echo.Context) (err error) {
	since := ctx.Get("since").(int64)
	limit := ctx.Get("limit").(int32)
	if limit == 0 {
		limit = 100
	}
	desc := ctx.Get("desc").(bool)

	var rows *pgx.Rows
	if since == 0 && !desc {
		rows, err = db.Read(0).Query("selectAllTeamsAsc", limit)
	} else if since == 0 && desc {
		rows, err = db.Read(0).Query("selectAllTeamsDesc", limit)
	} else if since != 0 && !desc {
		rows, err = db.Read(0).Query("selectAllTeamsAscSince", since, limit)
	} else {
		rows, err = db.Read(0).Query("selectAllTeamsDescSince", since, limit)
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

		team := models.Team{}
		err := rows.Scan(&team.ID, &team.Title, &team.Country, &team.Slug)
		if err != nil {
			ctx.Logger().Error(err)
		}

		blob, err := team.MarshalBinary()
		if err != nil {
			ctx.Logger().Error(err)
		}
		ctx.Response().Write(blob)
	}

	ctx.Response().Write([]byte("]"))
	ctx.Response().Flush()

	return nil
}

func GetTeam(ctx echo.Context) (err error) {
	slug := ctx.Get("slug").(string)
	shardIndex := ctx.Get("shardIndex").(int)

	team := models.Team{}
	err = db.Read(shardIndex).QueryRow("selectTeam", slug).Scan(&team.ID, &team.Title, &team.Country, &team.Slug)
	if err != nil {
		ctx.Logger().Error(err)

		apiError := models.Error{Message: "Not found"}
		blob, err := apiError.MarshalBinary()
		if err != nil {
			ctx.Logger().Error(err)
		}

		ctx.JSONBlob(http.StatusNotFound, blob)

		return err
	}

	blob, err := team.MarshalBinary()
	if err != nil {
		ctx.Logger().Error(err)
	}

	return ctx.JSONBlob(http.StatusOK, blob)
}

func AddTeam(ctx echo.Context) (err error) {
	bodyBlob, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		ctx.Logger().Error(err)
	}

	team := models.Team{}
	err = team.UnmarshalBinary(bodyBlob)
	if err != nil {
		ctx.Logger().Error(err)
	}

	shardIndex := int(adler32.Checksum(bodyBlob)) % db.GetShardsCount()
	tx, err := db.Write(shardIndex).Begin()
	if err != nil {
		ctx.Logger().Error(err)
	}

	team.Slug = fmt.Sprintf("%d-%s", shardIndex, slugify.Slugify(*team.Title))
	err = tx.QueryRow("insertTeam", team.Title, team.Country, team.Slug).Scan(&team.ID)
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

	blob, err := team.MarshalBinary()
	if err != nil {
		ctx.Logger().Error(err)
	}

	return ctx.JSONBlob(http.StatusCreated, blob)
}

func DeleteTeam(ctx echo.Context) (err error) {
	slug := ctx.Get("slug").(string)
	shardIndex := ctx.Get("shardIndex").(int)

	tx, err := db.Write(shardIndex).Begin()
	if err != nil {
		ctx.Logger().Error(err)
	}

	_, err = tx.Exec("deleteTeam", slug)
	if err != nil {
		ctx.Logger().Error(err)
		apiError := models.Error{Message: "Not found"}
		blob, err := apiError.MarshalBinary()
		if err != nil {
			ctx.Logger().Error(err)
		}

		err = tx.Rollback()
		if err != nil {
			ctx.Logger().Error(err)
		}

		ctx.JSONBlob(http.StatusNotFound, blob)

		return err
	}
	tx.Commit()

	return ctx.NoContent(http.StatusOK)
}
