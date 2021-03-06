package chess

import (
	"bytes"
	"fmt"
)

type Board struct {
	Squares [FILE_H + 1][RANK_8 + 1]Piecer
}

func InitialBoard() (b Board) {
	for _, t := range []int{TEAM_WHITE, TEAM_BLACK} {
		for f := FILE_A; f <= FILE_H; f++ {
			switch f {
			case FILE_A, FILE_H:
				r := &Rook{}
				r.Team = t
				b.Squares[f][StartRank(t)] = r
			case FILE_B, FILE_G:
				k := &Knight{}
				k.Team = t
				b.Squares[f][StartRank(t)] = k
			case FILE_C, FILE_F:
				bi := &Bishop{}
				bi.Team = t
				b.Squares[f][StartRank(t)] = bi
			case FILE_D:
				q := &Queen{}
				q.Team = t
				b.Squares[f][StartRank(t)] = q
			case FILE_E:
				k := &King{}
				k.Team = t
				b.Squares[f][StartRank(t)] = k
			}
			// Pawn
			p := &Pawn{}
			p.Team = t
			b.Squares[f][StartRank(t)+t] = p
		}
	}
	return
}

// Gets the starting rank of the team,
func StartRank(team int) int {
	if team == TEAM_WHITE {
		return RANK_1
	}
	return RANK_8
}

func EndRank(team int) int {
	return StartRank(-team)
}

// Whether a location is on the board
func IsValidLocation(l Location) bool {
	return l.Rank >= RANK_1 && l.Rank <= RANK_8 && l.File >= FILE_A &&
		l.File <= FILE_H
}

func (b Board) PieceAt(l Location) Piecer {
	return b.Squares[l.File][l.Rank]
}

func (b Board) IsEmpty(l Location) bool {
	return b.PieceAt(l) == nil
}

func (b Board) Render() string {
	buf := bytes.NewBuffer([]byte{' '})
	buf.WriteString(RenderRanks())
	buf.WriteByte('\n')
	for r := RANK_8; r >= RANK_1; r-- {
		buf.WriteString(fmt.Sprintf("%d ", r+1))
		for f := FILE_A; f <= FILE_H; f++ {
			l := Location{f, r}
			p := b.PieceAt(l)
			if p == nil {
				buf.WriteString(`{{c "gray"}}·{{_c}}`)
			} else {
				buf.WriteRune(p.Rune())
			}
			buf.WriteByte(' ')
		}
		buf.WriteString(fmt.Sprintf("%d\n", r+1))
	}
	buf.WriteByte(' ')
	buf.WriteString(RenderRanks())
	return buf.String()
}

func RenderRanks() string {
	buf := bytes.NewBuffer([]byte{})
	for _, r := range []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'} {
		buf.WriteString(` {{c "gray"}}`)
		buf.WriteByte(r)
		buf.WriteString("{{_c}}")
	}
	return buf.String()
}
