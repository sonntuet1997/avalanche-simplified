package constants

import "errors"

var (
	ErrRecordNotFound         = errors.New("record not found")
	ErrNotEnoughNeighborNodes = errors.New("not_enough_neighbor_nodes")
)
