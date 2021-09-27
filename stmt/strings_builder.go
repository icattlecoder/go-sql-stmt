package stmt

import (
	"strings"
	"sync"
)

var sbPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

func get() *strings.Builder {
	return sbPool.Get().(*strings.Builder)
}

func put(b *strings.Builder) {
	b.Reset()
	sbPool.Put(b)
}

func argsSqlString(sb *strings.Builder, args ...args) {
	for i, a := range args {
		sb.WriteString(a.SqlString())
		if i < len(args)-1 {
			sb.WriteString(", ")
		}
	}
}
