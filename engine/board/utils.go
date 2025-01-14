package board

func CoordToBoardIndex(coord string) byte {
	file := coord[0] - 'a'
	rank := int(coord[1]-'0') - 1

	return byte(rank*8) + file
}
