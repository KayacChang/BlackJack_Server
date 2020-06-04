package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gitlab.com/ti-backend/go-modules/ulgsdk/ulg168utils"
	"gitlab.com/ti-backend/ulg168/blackjack/game"
)

func main() {
	ulg168utils.InitConf(&ulg168utils.ULG168Conf{
		// APIHost: conf.ULG168APIHost,
		// MaintainAPI:   conf.MaintainAPI,
		// MaintainToken: conf.MaintainToken,
		// GameID:        conf.GameID,
		// GameType:      "blackjack",
		// ENV:           conf.Env,
		CMDPrefix: 8003000,
	})
	if err := game.Run(); err != nil {
		log.Println(err)
	}
}
