package keys

import (
	"errors"
	"time"
)

//ErrNoIDs no id is generated in queue
var ErrNoIDs = errors.New("uuid: no id in queue")

// IDGenerator uuid generator
type IDGenerator interface {
	NewID() (string, error)
	NewTimeID() (string, error)
	NewWith(t time.Time) (string, error)
}
