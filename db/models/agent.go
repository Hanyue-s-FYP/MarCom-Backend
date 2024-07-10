package models

import "database/sql"

type Agent struct {
	ID                 int
	Name               sql.NullString
	GeneralDescription sql.NullString
	BusinessID         int
	Attributes         []AgentAttribute
}

// no need embed agent here, just put all agent attributes relevant to the agent to the agent struct
type AgentAttribute struct {
	ID    int
	Key   sql.NullString
	Value sql.NullString
}
