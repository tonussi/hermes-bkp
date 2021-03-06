package kv

type Op uint64

const (
	GetOp  Op = 1
	SetOp  Op = 2
	DelOp  Op = 3
	SnapOp Op = 4
)
