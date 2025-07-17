package gpio

type Direction int

const (
	IN Direction = iota
	OUT
)

type Bank int

const (
	BANK_0 Bank = iota
	BANK_1
	BANK_2
	BANK_3
	BANK_4
)

type Group int

const (
	GROUP_A Group = iota
	GROUP_B
	GROUP_C
	GROUP_D
)

type X int

const (
	X_0 X = iota
	X_1
	X_2
	X_3
	X_4
	X_5
	X_6
	X_7
)

type OutLevel int

const (
	LOW OutLevel = iota
	HIGH
)
