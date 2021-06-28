package dom

type OperationType int

const(
	None = iota
	Eq
	In
	Lt
	Mod
)