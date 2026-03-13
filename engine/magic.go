package engine

import (
	"math/bits"
	"math/rand"
)

type MagicEntry struct {
	mask  Bitboard
	magic uint64
	shift uint
}

var RookMagics   [64]MagicEntry
var BishopMagics [64]MagicEntry
var RookAttacks  [64][]Bitboard
var BishopAttacks [64][]Bitboard

func rookOccMask(sq Square) Bitboard {
	return (HorizontalMasks[sq] & ^FileH & ^FileA) |
		(VerticalMasks[sq] & ^Rank1 & ^Rank8)
}

func bishopOccMask(sq Square) Bitboard {
	return (DiagonalMasks[sq] | AntiDiagonalMasks[sq]) &
		^Rank1 & ^Rank8 & ^FileH & ^FileA
}

func slidingAttacks(sq Square, occupied, mask Bitboard) Bitboard {
	squareBoard := SquareMasks[sq]
	bottom := ((occupied & mask) - (squareBoard << 1)) & mask
	top := ReverseBitboard(ReverseBitboard(occupied&mask) - 2*ReverseBitboard(squareBoard))
	return (bottom ^ top) & mask
}

func rookRefAttacks(sq Square, occ Bitboard) Bitboard {
	return slidingAttacks(sq, occ, HorizontalMasks[sq]) |
		slidingAttacks(sq, occ, VerticalMasks[sq])
}

func bishopRefAttacks(sq Square, occ Bitboard) Bitboard {
	return slidingAttacks(sq, occ, DiagonalMasks[sq]) |
		slidingAttacks(sq, occ, AntiDiagonalMasks[sq])
}

func InitMagics() {
	rng := rand.New(rand.NewSource(12345))

	for sq := H1; sq <= A8; sq++ {
		// Rook
		{
			mask := rookOccMask(sq)
			n := bits.OnesCount64(uint64(mask))
			shift := uint(64 - n)
			size := 1 << n

			occs := make([]Bitboard, size)
			atts := make([]Bitboard, size)

			occ := Bitboard(0)
			for i := 0; i < size; i++ {
				occs[i] = occ
				atts[i] = rookRefAttacks(sq, occ)
				occ = (occ - mask) & mask
			}

			table := make([]Bitboard, size)
			for {
				candidate := rng.Uint64() & rng.Uint64() & rng.Uint64()
				for i := range table {
					table[i] = 0
				}
				fail := false
				for i := 0; i < size; i++ {
					idx := (occs[i] * Bitboard(candidate)) >> shift
					if table[idx] == 0 {
						table[idx] = atts[i]
					} else if table[idx] != atts[i] {
						fail = true
						break
					}
				}
				if !fail {
					RookMagics[sq] = MagicEntry{mask: mask, magic: candidate, shift: shift}
					RookAttacks[sq] = make([]Bitboard, size)
					copy(RookAttacks[sq], table)
					break
				}
			}
		}

		// Bishop
		{
			mask := bishopOccMask(sq)
			n := bits.OnesCount64(uint64(mask))
			shift := uint(64 - n)
			size := 1 << n

			occs := make([]Bitboard, size)
			atts := make([]Bitboard, size)

			occ := Bitboard(0)
			for i := 0; i < size; i++ {
				occs[i] = occ
				atts[i] = bishopRefAttacks(sq, occ)
				occ = (occ - mask) & mask
			}

			table := make([]Bitboard, size)
			for {
				candidate := rng.Uint64() & rng.Uint64() & rng.Uint64()
				for i := range table {
					table[i] = 0
				}
				fail := false
				for i := 0; i < size; i++ {
					idx := (occs[i] * Bitboard(candidate)) >> shift
					if table[idx] == 0 {
						table[idx] = atts[i]
					} else if table[idx] != atts[i] {
						fail = true
						break
					}
				}
				if !fail {
					BishopMagics[sq] = MagicEntry{mask: mask, magic: candidate, shift: shift}
					BishopAttacks[sq] = make([]Bitboard, size)
					copy(BishopAttacks[sq], table)
					break
				}
			}
		}
	}
}
