package handlers

import (
	"github.com/labstack/echo"
	"io/ioutil"
	"github.com/oleggator/esports-backend/models"
	"github.com/oleggator/esports-backend/utils"
	"net/http"
	"github.com/oleggator/esports-backend/db"
)

func GetPlayers(ctx echo.Context) (err error) {
	return nil
}

func GetPlayer(ctx echo.Context) (err error) {
	return nil
}

func AddPlayer(ctx echo.Context) (err error) {
	bodyBlob, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		ctx.Logger().Error(err)
	}

	player := models.Player{}
	err = player.UnmarshalBinary(bodyBlob)
	if err != nil {
		ctx.Logger().Error(err)
	}
	if player.Team == nil {
		apiError := models.Error{Message: "Team field is empty"}
		blob, err := apiError.MarshalBinary()
		if err != nil {
			ctx.Logger().Error(err)
		}

		ctx.JSONBlob(http.StatusNotFound, blob)
	}

	shardIndex, err := utils.GetShardIndex(player.Team)
	if err != nil {
		ctx.Logger().Error(err)

		apiError := models.Error{Message: "Team not found"}
		blob, err := apiError.MarshalBinary()
		if err != nil {
			ctx.Logger().Error(err)
		}

		ctx.JSONBlob(http.StatusNotFound, blob)
	}

	var teamId int64
	err = db.Read(shardIndex).QueryRow("getTeamBySlug", player.Team).Scan(&teamId, &player.Team)
	if err != nil {
		ctx.Logger().Error(err)

		apiError := models.Error{Message: "Team not found"}
		blob, err := apiError.MarshalBinary()
		if err != nil {
			ctx.Logger().Error(err)
		}

		ctx.JSONBlob(http.StatusNotFound, blob)
	}

	tx, err := db.Write(shardIndex).Begin()
	if err != nil {
		ctx.Logger().Error(err)
	}

	err = tx.QueryRow("insertPlayer",
		player.Fullname, player.Nickname, player.Country, player.Team).Scan(&player.ID)
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

	tx.Commit()
	if err != nil {
		ctx.Logger().Error(err)
	}

	blob, err := player.MarshalBinary()
	if err != nil {
		ctx.Logger().Error(err)
	}

	return ctx.JSONBlob(http.StatusCreated, blob)
}

func DeletePlayer(ctx echo.Context) (err error) {
	return nil
}
