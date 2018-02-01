package main

import (
	"github.com/labstack/echo"
	"strconv"
	"net/http"
	"github.com/oleggator/esports-backend/models"
	"strings"
	"github.com/oleggator/esports-backend/db"
)

func ArrayMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		if sinceString := ctx.QueryParam("since"); sinceString != "" {
			value, err := strconv.ParseInt(sinceString, 0, 64)
			if err != nil {
				ctx.Logger().Error(err)
				apiErr := models.Error{"Wrong since parameter"}
				blob, err := apiErr.MarshalBinary()
				if err != nil {
					ctx.Logger().Error(err)
				}

				return ctx.JSONBlob(http.StatusBadRequest, blob)
			}

			ctx.Set("since", value)
		} else {
			ctx.Set("since", int64(0))
		}

		if descString := ctx.QueryParam("desc"); descString != "" {
			value, err := strconv.ParseBool(descString)
			if err != nil {
				ctx.Logger().Error(err)

				apiErr := models.Error{"Wrong desc parameter"}
				blob, err := apiErr.MarshalBinary()
				if err != nil {
					ctx.Logger().Error(err)
				}

				return ctx.JSONBlob(http.StatusBadRequest, blob)
			}

			ctx.Set("desc", value)
		} else {
			ctx.Set("desc", false)
		}

		if limitString := ctx.QueryParam("limit"); limitString != "" {
			value, err := strconv.ParseInt(limitString, 0, 32)
			if err != nil {
				ctx.Logger().Error(err)

				apiErr := models.Error{"Wrong limit parameter"}
				blob, err := apiErr.MarshalBinary()
				if err != nil {
					ctx.Logger().Error(err)
				}

				return ctx.JSONBlob(http.StatusBadRequest, blob)
			}

			ctx.Set("limit", int32(value))
		} else {
			ctx.Set("limit", int32(0))
		}

		return next(ctx)
	}
}

func GetByIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if idString := ctx.Param("id"); idString != "" {
			id, err := strconv.ParseInt(idString, 0, 64)
			if err != nil {
				apiErr := models.Error{"Wrong id"}
				blob, err := apiErr.MarshalBinary()
				if err != nil {
					ctx.Logger().Error(err)
				}

				return ctx.JSONBlob(http.StatusBadRequest, blob)
			}

			ctx.Set("id", int64(id))
			return next(ctx)
		} else {
			apiErr := models.Error{"Empty id parameter"}
			blob, err := apiErr.MarshalBinary()
			if err != nil {
				ctx.Logger().Error(err)
			}

			return ctx.JSONBlob(http.StatusBadRequest, blob)
		}
	}
}

func GetBySlugMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if slugString := ctx.Param("slug"); slugString != "" {
			i := strings.IndexRune(slugString, '-')
			if i == -1 {
				apiErr := models.Error{"Wrong slug"}
				blob, err := apiErr.MarshalBinary()
				if err != nil {
					ctx.Logger().Error(err)
				}

				return ctx.JSONBlob(http.StatusBadRequest, blob)
			}

			shardIndex, err := strconv.Atoi(slugString[:i])
			if err != nil {
				apiErr := models.Error{"Wrong slug"}
				blob, err := apiErr.MarshalBinary()
				if err != nil {
					ctx.Logger().Error(err)
				}

				return ctx.JSONBlob(http.StatusBadRequest, blob)
			}

			if shardIndex >= db.GetShardsCount() {
				apiErr := models.Error{Message: "Not found"}
				blob, err := apiErr.MarshalBinary()
				if err != nil {
					ctx.Logger().Error(err)
				}

				return ctx.JSONBlob(http.StatusNotFound, blob)
			}

			ctx.Set("slug", slugString)
			ctx.Set("shardIndex", shardIndex)

			return next(ctx)
		} else {
			apiErr := models.Error{"Empty slug parameter"}
			blob, err := apiErr.MarshalBinary()
			if err != nil {
				ctx.Logger().Error(err)
			}

			return ctx.JSONBlob(http.StatusBadRequest, blob)
		}
	}
}

func GetByNicknameMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if nicknameString := ctx.Param("nickname"); nicknameString != "" {
			ctx.Set("nickname", nicknameString)
			return next(ctx)
		} else {
			apiErr := models.Error{"Empty nickname parameter"}
			blob, err := apiErr.MarshalBinary()
			if err != nil {
				ctx.Logger().Error(err)
			}

			return ctx.JSONBlob(http.StatusBadRequest, blob)
		}
	}
}