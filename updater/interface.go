package updater

import (
	"time"
)

type Updater interface {
	Update(time.Duration, string)
	Name() string
}
