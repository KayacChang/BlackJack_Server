package table

import (
	"database/sql"

	"gitlab.com/ti-backend/ulg168/blackjack/protocol"
	"gitlab.com/ti-backend/go-modules/casino/lobby"
	"go.uber.org/zap"
)

type Options struct {
	protocol.TableStat
	Database *sql.DB
	Logger   *zap.Logger
	Lobby    *lobby.Lobby
	LastShoe int
}
