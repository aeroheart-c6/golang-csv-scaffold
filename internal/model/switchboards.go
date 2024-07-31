package model

import "time"

type Switchboard struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SwitchboardPanel struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
