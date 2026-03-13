package engine

import (
	"math/bits"
	"math/rand"
)

type MagicEntry struct {
	mask   Bitboard
	magic  uint64
	shift  uint
	offset uint32
}

var RookMagics [64]MagicEntry
var BishopMagics [64]MagicEntry

var rookTable []Bitboard
var bishopTable []Bitboard

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

	// First pass: compute total table sizes
	rookTotalSize := 0
	bishopTotalSize := 0

	for sq := H1; sq <= A8; sq++ {
		rookN := bits.OnesCount64(uint64(rookOccMask(sq)))
		rookTotalSize += 1 << rookN

		bishopN := bits.OnesCount64(uint64(bishopOccMask(sq)))
		bishopTotalSize += 1 << bishopN
	}

	// Allocate flat tables
	rookTable = make([]Bitboard, rookTotalSize)
	bishopTable = make([]Bitboard, bishopTotalSize)

	rookOffset := uint32(0)
	bishopOffset := uint32(0)

	for sq := H1; sq <= A8; sq++ {
		// === Rook ===
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

			// Find magic (same trial-and-error loop as before)
			for {
				candidate := rng.Uint64() & rng.Uint64() & rng.Uint64()

				// Clear the region of rookTable we're about to write
				for i := uint32(0); i < uint32(size); i++ {
					rookTable[rookOffset+i] = 0
				}

				fail := false
				for i := 0; i < size; i++ {
					idx := uint32((occs[i] * Bitboard(candidate)) >> shift)
					entry := rookOffset + idx
					if rookTable[entry] == 0 {
						rookTable[entry] = atts[i]
					} else if rookTable[entry] != atts[i] {
						fail = true
						break
					}
				}
				if !fail {
					RookMagics[sq] = MagicEntry{
						mask:   mask,
						magic:  candidate,
						shift:  shift,
						offset: rookOffset,
					}
					rookOffset += uint32(size)
					break
				}
			}
		}

		// === Bishop ===
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

			for {
				candidate := rng.Uint64() & rng.Uint64() & rng.Uint64()

				for i := uint32(0); i < uint32(size); i++ {
					bishopTable[bishopOffset+i] = 0
				}

				fail := false
				for i := 0; i < size; i++ {
					idx := uint32((occs[i] * Bitboard(candidate)) >> shift)
					entry := bishopOffset + idx
					if bishopTable[entry] == 0 {
						bishopTable[entry] = atts[i]
					} else if bishopTable[entry] != atts[i] {
						fail = true
						break
					}
				}
				if !fail {
					BishopMagics[sq] = MagicEntry{
						mask:   mask,
						magic:  candidate,
						shift:  shift,
						offset: bishopOffset,
					}
					bishopOffset += uint32(size)
					break
				}
			}
		}
	}
}
