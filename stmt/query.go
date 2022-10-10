package stmt

import (
	"github.com/keegancsmith/sqlf"
)

type Querier interface {
	Query() *sqlf.Query
}
