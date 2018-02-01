package main

import (
	"github.com/oleggator/esports-backend/db"
	"github.com/oleggator/esports-backend/handlers"
	"github.com/labstack/echo"
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

func main() {
	config := getConfig("config.yml")
	err := db.InitDB(config.DB)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()


	e := echo.New()

	e.GET("/games", handlers.GetGames, ArrayMiddleware)
	e.GET("/games/:slug", handlers.GetGame, GetBySlugMiddleware)
	e.POST("/games", handlers.AddGame)
	e.DELETE("/games/:slug", handlers.DeleteGame, GetBySlugMiddleware)

	e.GET("/teams", handlers.GetTeams, ArrayMiddleware)
	e.GET("/teams/:slug", handlers.GetTeam, GetBySlugMiddleware)
	e.POST("/teams", handlers.AddTeam)
	e.DELETE("/teams/:slug", handlers.DeleteTeam, GetBySlugMiddleware)

	e.GET("/teams/:slug/players", handlers.GetPlayers, GetBySlugMiddleware, ArrayMiddleware)
	e.GET("/teams/:slug/players/:nickname", handlers.GetPlayer, GetBySlugMiddleware, GetByNicknameMiddleware)
	e.POST("/teams/:slug/players", handlers.AddPlayer, GetBySlugMiddleware)
	e.DELETE("/teams/:slug/players", handlers.DeletePlayer, GetBySlugMiddleware, )

	//e.GET("/players", handlers.GetPlayers, ArrayMiddleware)
	//e.GET("/players/:nickname", handlers.GetPlayer, GetByNicknameMiddleware)
	//e.POST("/players", handlers.AddPlayer)
	//e.DELETE("/players/:nickname", handlers.GetPlayers)

	e.Logger.Fatal(e.Start(config.HTTP.Address))
}

func getConfig(path string) (config Config) {
	configBlob, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("Config error:", err)
	}

	err = yaml.Unmarshal(configBlob, &config)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}

	return config
}
