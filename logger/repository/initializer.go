package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mailru/go-clickhouse"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

func CreateDatabase(dbURL, configDir string) *sqlx.DB {
	log.Print("Connect to db url ", dbURL, " ...")
	c, err := sqlx.Connect("clickhouse", dbURL)
	if err != nil {
		log.Print(err)
		time.Sleep(time.Second * 15)
		CreateDatabase(dbURL, configDir)
	}
	log.Print("Connected to db.", " Init schema...")

	if configDir != "" {
		ddl, err := ioutil.ReadFile(fmt.Sprintf("%s/config/schema.sql", configDir))
		if err != nil {
			log.Fatal(err)
		}
		log.Print(ddl)
		log.Print("Read schema from file...")
		if _, err := c.Exec(string(ddl)); err != nil {
			if strings.Contains(err.Error(), "Code: 57") {
				newddl := strings.ReplaceAll(string(ddl), "create table if not exists", "ATTACH TABLE")
				if _, err := c.Exec(newddl); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}
		log.Print("Schema created.")
	}

	return c
}
