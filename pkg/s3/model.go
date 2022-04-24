package s3

import "time"

type Object struct {
	Key        string
	Size       int64
	ModifiedAt time.Time
}
