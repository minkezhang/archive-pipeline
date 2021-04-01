package interfaces

import (
	"github.com/minkezhang/archive-pipeline/import/record"
)

type I interface {
	Import() (*record.R, error)
}
