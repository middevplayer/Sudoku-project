package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/eiannone/keyboard"
)

// Создаем наши константы, которые будем использовать на протяжении всей игры
const (
	rows    = 9 // длина по y
	columns = 9 // длина по x
	empty   = 0 // пустая клетка
)

type Cell struct { //создаем нашу клетку, которая принимает цифры и закрепляет их, если они были изначально
	digit int8
	fix   bool
}

type Grid [rows][columns]Cell //создаем нашу сетку то есть двухмерный массив 9 на 9 заполненный нашими клетками

func (g *Grid) InBounds(row, column int) bool { //функция на проверку границ
	if row < 0 || row >= rows || column < 0 || column >= columns {
		return false
	}
	return true

}

var ( //создаем наши ошибки для вывода если что то пошло не так
	ErrBounds = errors.New("Вне поля")
	ErrFixed  = errors.New("Эту цифру нельзя менять")
	ErrValid  = errors.New("Неправильная цифра (1-9)")
	ErrRow    = errors.New("Эта цифра уже есть в ряду")
	ErrColumn = errors.New("Эта цифра уже есть в столбе")
	ErrRegion = errors.New("Эта цифра уже есть в этом квадрате")
)

func (g *Grid) IsFixed(row, column int) bool { //функция-проверка отвечающая за замену изначального числа
	return g[row][column].fix
}
func (g *Grid) Set(row, column int, digit int8) error { //функция для проверки всей доски и последующего постановки цифры на нее
	switch {
	case !g.InBounds(row, column):
		return ErrBounds
	case g.IsFixed(row, column):
		return ErrFixed
	// Проверяем команду на стирание ДО проверки ValidDigits
	case digit == empty:
		g[row][column].digit = empty
		return nil
	case !ValidDigits(digit):
		return ErrValid
	case g.InColumn(column, digit):
		return ErrColumn
	case g.InRow(row, digit):
		return ErrRow
	case g.InRegion(row, column, digit):
		return ErrRegion
	}
	g[row][column].digit = digit //цифра поставилась если прошла все проверки
	return nil
}

func ValidDigits(digit int8) bool { //валидность числа
	return digit >= 1 && digit <= 9
}

func (g *Grid) InRow(row int, digit int8) bool { //проверка на то что в строке нет одинаковый цифр
	for c := 0; c < columns; c++ {
		if g[row][c].digit == digit {
			return true
		}
	}
	return false
}

func (g *Grid) InColumn(column int, digit int8) bool { //проверка на то что в столбе нет одинаковый цифр
	for r := 0; r < rows; r++ {
		if g[r][column].digit == digit {
			return true
		}
	}
	return false
}

func (g *Grid) InRegion(row, column int, digit int8) bool { // проверка на разные цифры в 3 на 3 квадрате
	startRow := row / 3 * 3       // делаем начальные позиции. Из за того что мы работаем в инт он отбрасывает после запятой, так что 4/3 * 3 = 3. Таким образом делим все на 3 квадрата
	startColumn := column / 3 * 3 // тут также

	for r := startRow; r < startRow+3; r++ { // перебор маленького квадрата 3 на 3. Работает также как и с большим квадратом
		for c := startColumn; c < startColumn+3; c++ {
			if g[r][c].digit == digit {
				return true
			}
		}
	}
	return false
}
func (g *Grid) SudokugenerateBot(initalNumbersCount int) bool {
	for i := 0; i < initalNumbersCount; {
		x := rand.IntN(9)
		y := rand.IntN(9)
		value := int8(rand.IntN(9) + 1)

		// 1. Если клетка уже занята, ищем другую
		if g[x][y].digit != empty {
			continue
		}

		// 2. Проверяем, свободна ли строка, столбец и квадрат от этой цифры
		if !g.InRow(x, value) && !g.InColumn(y, value) && !g.InRegion(x, y, value) {
			g[x][y].digit = value // Ставим, если всё чисто
			g[x][y].fix = true
			i++ // Переходим к следующей цифре
		}
	}

	return true
}

func (g *Grid) SudokuBot() bool { // автобот работающий по перебору. Ставит в пустую клетку число от 1 до 9 и потом решает основываясь на этой цифре.
	//  Если случится так, что доска не подходит под правила, то ставит другую цифру и так далее
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

func (g *Grid) Draw(curRow, curCol int) { //функция для вывода доски
	// Глубокий темно-черничный фон для всего поля
	bgMain := "\033[48;2;16;24;48m"
	// Чуть более светлый синий фон для курсора, чтобы он выделялся
	bgCursor := "\033[48;2;35;55;95m"

	// сине-серый цвет для разделительных линий
	textGrid := "\033[38;2;70;85;110m"

	// СВЕТЛО-ГОЛУБОЙ цвет для (зафиксированных) цифр
	textFixedLight := "\033[38;2;160;210;255m"

	// СВЕТЛО-РОЗОВЫЙ ЖИРНЫЙ
	textUserPink := "\033[1m\033[38;2;255;165;210m"

	reset := "\033[0m" // Полный сброс всех стилей

	for i := 0; i < rows; i++ {
		if i%3 == 0 && i != 0 {
			fmt.Print(bgMain + textGrid + "------+-------+------ " + reset + "\n")
		}

		for j := 0; j < columns; j++ {
			if j%3 == 0 && j != 0 {
				fmt.Print(bgMain + textGrid + "| " + reset)
			}

			if i == curRow && j == curCol {
				// КЛЕТКА КУРСОРA
				if g[i][j].digit == empty {
					fmt.Print(bgCursor + textGrid + ". " + reset)
				} else {
					if g[i][j].fix {
						fmt.Printf(bgCursor+textFixedLight+"%d "+reset, g[i][j].digit)
					} else {
						fmt.Printf(bgCursor+textUserPink+"%d "+reset, g[i][j].digit)
					}
				}
			} else {
				// ОБЫЧНЫЕ КЛЕТКИ (
				if g[i][j].digit == empty {
					fmt.Print(bgMain + textGrid + ". " + reset)
				} else {
					if g[i][j].fix {
						// Стартовые цифры (Светло-голубые)
						fmt.Printf(bgMain+textFixedLight+"%d "+reset, g[i][j].digit)
					} else {
						// Ваши цифры (Светло-розовые)
						fmt.Printf(bgMain+textUserPink+"%d "+reset, g[i][j].digit)
					}
				}
			}
		}
		fmt.Println()
	}
}

func (g *Grid) SudokuParser(str string) { //так как тут работает все на шаблонах. Для удобства я написал функцию для вывода шаблона в матрицу из строки.
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

func (g *Grid) GridSelect() { //тут пишем свои судоку в строку. Точка = пустота. Также это выбор сложности
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
	case 5:
		g.SudokugenerateBot(20)

	}
}
func (g *Grid) IsSolved() bool { //легкая проверка доски на решение. Проверяет
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
func (g *Grid) Move(cursorRow, cursorCol int) (int, int, int8) {
	char, key, err := keyboard.GetSingleKey()
	if err != nil {
		panic(err)
	}
	if key == keyboard.KeyEsc {
		fmt.Println("\nВыход из игры")
		os.Exit(0)
	}
	if char == 0 {
		switch key {
		case keyboard.KeyArrowUp:
			if cursorRow > 0 {
				cursorRow--
			}
		case keyboard.KeyArrowDown:
			if cursorRow < rows-1 {
				cursorRow++
			}
		case keyboard.KeyArrowLeft:
			if cursorCol > 0 {
				cursorCol--
			}
		case keyboard.KeyArrowRight:
			if cursorCol < columns-1 {
				cursorCol++
			}
		case keyboard.KeyBackspace, keyboard.KeyBackspace2:
			return cursorRow, cursorCol, 0 // Сигнал стереть цифру
		case keyboard.KeySpace:
			return cursorRow, cursorCol, 99 // Сигнал для бота-авторешения
		}
		return cursorRow, cursorCol, -1 // Просто подвинули курсор, цифру не вводили
	}
	if char >= '0' && char <= '9' {
		return cursorRow, cursorCol, int8(char - '0')
	}
	return cursorRow, cursorCol, -1
}

func (g *Grid) Theend() { // когда сетка подходит под функцию IsSolved тогда ты победил.
	if !g.IsSolved() {
		return
	}

	fmt.Println("Ты победил!")
	os.Exit(0)
}
func main() { //основная функция
	fmt.Println("Добро пожаловать в Судоку на курсе Go")
	var board Grid
	fmt.Println("Выберите сложность доски: 1 - легкая, 2 - средняя, 3 - сложная, 4 - мастер, 5 рандом")
	board.GridSelect()
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	cursorRow, cursorCol := 0, 0 //стартовые позиции
	var lastErr error

	for {
		fmt.Print("\033[H\033[2J") // чтобы стереть предыдущий кадр
		fmt.Println("Управление: СТРЕЛКИ - движение | 1-9 - ввод цифры | 0/Backspace - стереть | Пробел - автобот | Esc - выход")
		if lastErr != nil {
			fmt.Printf("Ошибка %v\n", lastErr)
			lastErr = nil
		} else {
			fmt.Println()
		}
		board.Draw(cursorRow, cursorCol)
		var digit int8
		cursorRow, cursorCol, digit = board.Move(cursorRow, cursorCol)
		if digit == 99 {
			if board.SudokuBot() {
				fmt.Print("\033[H\033[2J")
				fmt.Println("Судоку успешно решено ботом:")
				board.Draw(-1, -1) // Рисуем без подсветки курсора
				break
			} else {
				lastErr = errors.New("Решения для этой доски не существует")
			}
			continue
		}
		if digit != -1 {
			err := board.Set(cursorRow, cursorCol, digit)
			if err != nil {
				lastErr = err
			} else {
				board.Theend()
			}
		}
	}
}
