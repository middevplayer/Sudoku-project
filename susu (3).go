package main

import (
	"errors"
	"fmt"
	"os"
)

const (
	rows       = 9
	columns    = 9
	empty      = 0
	colorGreen = "\033[32m"
	ColorReset = "\033[0m"
)

type Cell struct {
	digit int8
	fix   bool
}

type Grid [rows][columns]Cell

func (g *Grid) InBounds(row, column int) bool {
	if row < 0 || row >= rows || column < 0 || column >= columns {
		return false
	}
	return true

}

var (
	ErrBounds = errors.New("Вне поля")
	ErrFixed  = errors.New("Эту цифру нельзя менять")
	ErrValid  = errors.New("Неправильная цифра (1-9)")
	ErrRow    = errors.New("Эта цифра уже есть в ряду")
	ErrColumn = errors.New("Эта цифра уже есть в столбе")
	ErrRegion = errors.New("Эта цифра уже есть в этом квадрате")
)

func (g *Grid) IsFixed(row, column int) bool {
	return g[row][column].fix
}
func (g *Grid) Set(row, column int, digit int8) error {
	switch {
	case !g.InBounds(row, column):
		return ErrBounds
	case !ValidDigits(digit):
		return ErrValid
	case g.IsFixed(row, column):
		return ErrFixed

	case g.InColumn(column, digit):
		return ErrColumn

	case g.InRow(row, digit):
		return ErrRow

	case g.InRegion(row, column, digit):
		return ErrRegion
	}
	g[row][column].digit = digit
	return nil

}

func ValidDigits(digit int8) bool {
	return digit >= 1 && digit <= 9
}

func (g *Grid) InRow(row int, digit int8) bool {
	for c := 0; c < columns; c++ {
		if g[row][c].digit == digit {
			return true
		}
	}
	return false
}

func (g *Grid) InColumn(column int, digit int8) bool {
	for r := 0; r < rows; r++ {
		if g[r][column].digit == digit {
			return true
		}
	}
	return false
}

func (g *Grid) InRegion(row, column int, digit int8) bool {
	startRow := row / 3 * 3
	startColumn := column / 3 * 3

	for r := startRow; r < startRow+3; r++ {
		for c := startColumn; c < startColumn+3; c++ {
			if g[r][c].digit == digit {
				return true
			}
		}
	}
	return false
}

func (g *Grid) SudokuBot() bool {
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			if g[i][j].digit == empty {
				for d := int8(1); d <= 9; d++ {
					if !g.InRow(i, d) && !g.InColumn(j, d) && !g.InRegion(i, j, d) {
						g[i][j].digit = d
						if g.SudokuBot() {
							return true
						}
						g[i][j].digit = empty
					}
				}
				return false
			}
		}
	}
	return true
}

func (g *Grid) Draw() {
	var stickPrinted = false
	var horizontstick = false
	for i := 0; i < rows; i++ {
		if i%3 == 0 && horizontstick == false && i != 0 {
			fmt.Println()
			i--
			horizontstick = true
			continue
		}
		horizontstick = false
		for j := 0; j < columns; j++ {
			if j%3 == 0 && stickPrinted == false && j != 0 {
				fmt.Print("| ")
				j--
				stickPrinted = true
				continue
			}
			stickPrinted = false
			if g[i][j].digit == empty {
				fmt.Print(". ")
			} else {
				if g[i][j].fix {
					fmt.Printf("%d ", g[i][j].digit)
				} else {
					fmt.Printf("%s%d%s ", colorGreen, g[i][j].digit, ColorReset)
				}
			}

		}
		fmt.Println()

	}

}
func (g *Grid) SudokuParser(str string) {
	for i := 0; i < 81; i++ {
		r := i / 9
		c := i % 9
		if str[i] == '.' {
			g[r][c].digit = empty
			g[r][c].fix = false
		} else {
			g[r][c].digit = int8(str[i] - '0')
			g[r][c].fix = true
		}
	}
}

func (g *Grid) GridSelect() {
	ea := "82..4..6...16..89...98315.749.157.............53..4...96.415..81..7632..3...28.51"
	hd := ".......14......2.38...5.......2.7....31............65.6.....7.....14.......3....."
	md := "2.34....58.916.7.4..6.3..197.2..3.6...825......16.7..2..7..592693.72....6...9.47."
	ms := "3...49......6..5.1752..1.....1...7..5..396.....815..96..3.1..6...4...1......28..."
	var v int8
	fmt.Scan(&v)
	switch v {
	case 1:
		g.SudokuParser(ea)
	case 2:
		g.SudokuParser(md)
	case 3:
		g.SudokuParser(hd)
	case 4:
		g.SudokuParser(ms)

	}
}
func (g *Grid) IsSolved() bool {
	for r := 0; r < rows; r++ {
		for c := 0; c < columns; c++ {
			if g[r][c].digit == empty {
				return false
			}
		}
	}

	for r := 0; r < rows; r++ {
		for d := int8(1); d <= 9; d++ {
			if !g.InRow(r, d) {
				return false
			}
		}
	}

	for c := 0; c < columns; c++ {
		for d := int8(1); d <= 9; d++ {
			if !g.InColumn(c, d) {
				return false
			}
		}
	}

	return true
}

func (g *Grid) Theend() {
	if !g.IsSolved() {
		return
	}

	fmt.Println("Ты победил!")
	os.Exit(0)
}
func main() {
	fmt.Println("Добро пожаловать в Судоку на курсе Go")
	var board Grid
	fmt.Println("Выберите сложность доски: 1 - легкая, 2 - средняя, 3 - сложная, 4 - мастер")
	board.GridSelect()
	var row, column int
	var digit int8
	for {
		fmt.Print("\033[H\033[2J")
		fmt.Println("\nТекущая доска:")
		board.Draw()
		fmt.Println("Введите через пробел: столб(1-9)(x координата) строчка(1-9)(y координата) цифра(1-9) или 0 0 0 для авторешения")
		fmt.Scan(&column, &row, &digit)
		if column == 0 && row == 0 && digit == 0 {
			if board.SudokuBot() {
				fmt.Println("Судоку решено ботом")
			} else {
				fmt.Println("Решения не существует")
			}
			board.Draw()
			break
		}

		err := board.Set(row-1, column-1, digit)
		if err != nil {
			fmt.Printf("Ошибка холла: %v\n", err)
		} else {
			fmt.Println("Ход сделан")

			board.Theend()
		}

	}

}
