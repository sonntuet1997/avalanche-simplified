package entities

import "sync"

type TransactionNode struct {
	Transaction          *Transaction `json:"transaction"`
	Chit                 bool         `json:"chit"`
	Confidence           uint         `json:"confidence"`
	ConsecutiveSuccesses uint         `json:"consecutive_successes"`
	ParentID             string       `json:"parent_id"`
	ParentNode           *Transaction
	Mutex                sync.RWMutex
	ChildNodes           []*Transaction
}
