package chess

import (
	"bytes"
)

type Board struct {
	Squares [FILE_H + 1][RANK_8 + 1]Piecer
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
	buf := bytes.NewBuffer([]byte{})
	for r := RANK_8; r >= RANK_1; r-- {
		for f := FILE_A; f <= FILE_H; f++ {
			l := Location{f, r}
			p := b.PieceAt(l)
			if p == nil {
				buf.WriteString(`{{c "gray"}}•{{_c}}`)
			} else {
				buf.WriteRune(p.Rune())
			}
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}
