package db

import (
	"github.com/jackc/pgx"
	"log"
)

func prepareGames(conn *pgx.Conn) {
	var err error

	// Select games
	_, err = conn.Prepare("selectAllGamesAsc",
		"select id, title, slug from games order by title asc limit $1")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.Prepare("selectAllGamesDesc",
		"select id, title, slug from games order by title desc limit $1")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.Prepare("selectAllGamesAscSince",
		"select id, title, slug from games where id > $1 order by title asc limit $2")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.Prepare("selectAllGamesDescSince",
		"select id, title, slug from games where id < $1 order by title desc limit $2")
	if err != nil {
		log.Fatalln(err)
	}


	// Select game
	_, err = conn.Prepare("selectGame",
		"select id, title, slug from games where slug = $1 limit 1")
	if err != nil {
		log.Fatalln(err)
	}


	// Add game
	_, err = conn.Prepare("insertGame",
		"insert into games (title, slug) values ($1, $2) returning id")
	if err != nil {
		log.Fatalln(err)
	}


	// Delete game
	_, err = conn.Prepare("deleteGame",
		"delete from games where slug = $1")
	if err != nil {
		log.Fatalln(err)
	}
}

func prepareTeams(conn *pgx.Conn) {
	var err error

	// Select teams
	_, err = conn.Prepare("selectAllTeamsAsc",
		"select id, title, country, slug from teams order by title asc limit $1")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.Prepare("selectAllTeamsDesc",
		"select id, title, country, slug from teams order by title desc limit $1")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.Prepare("selectAllTeamsAscSince",
		"select id, title, country, slug from teams where id > $1 order by title asc limit $2")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.Prepare("selectAllTeamsDescSince",
		"select id, title, country, slug from teams where id < $1 order by title desc limit $2")
	if err != nil {
		log.Fatalln(err)
	}


	// Select team
	_, err = conn.Prepare("selectTeam",
		"select id, title, country, slug from teams where slug = $1 limit 1")
	if err != nil {
		log.Fatalln(err)
	}


	// Add team
	_, err = conn.Prepare("insertTeam",
		"insert into teams (title, country, slug) values ($1, $2, $3) returning id")
	if err != nil {
		log.Fatalln(err)
	}


	// Delete team
	_, err = conn.Prepare("deleteTeam",
		"delete from teams where slug = $1")
	if err != nil {
		log.Fatalln(err)
	}
}

func preparePlayers(conn *pgx.Conn) {
	var err error

	// Add player
	_, err = conn.Prepare("insertPlayer",
		"insert into players (fullname, nickname, country, teamSlug, team) values ($1, $2, $3, $4, $5) returning id")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.Prepare("getTeamBySlug",
		"select id, slug from teams where slug = $1 limit 1")
	if err != nil {
		log.Fatalln(err)
	}
}
