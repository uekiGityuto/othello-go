package model

import (
	"fmt"
)

type Color int

const (
	Black Color = iota
	White
)

func (c Color) Valid() error {
	if c == Black || c == White {
		return nil
	} else {
		return fmt.Errorf(`"%d" is invalid as color`, c)
	}
}

func (c Color) GetReversed() Color {
	if c == Black {
		return White
	} else {
		return Black
	}
}

func (c Color) String() string {
	if c == Black {
		return "黒"
	} else {
		return "白"
	}
}

type Stone struct {
	color Color
}

func NewStone(color Color) (*Stone, error) {
	if err := color.Valid(); err != nil {
		return nil, fmt.Errorf("failed to create stone: %w", err)
	}
	return &Stone{color: color}, nil
}

func (s *Stone) Get() string {
	if s.color == White {
		return "o"
	} else {
		return "●"
	}
}

func (s *Stone) IsBlack() bool {
	if s.color == Black {
		return true
	} else {
		return false
	}
}

func (s *Stone) Reverse() {
	s.color = s.color.GetReversed()
}

type Cell struct {
	stone *Stone
}

func NewCell(color *Color) (*Cell, error) {
	if color == nil {
		return &Cell{stone: nil}, nil
	}
	stone, err := NewStone(*color)
	if err != nil {
		return nil, fmt.Errorf("failed to create cell: %w", err)
	}
	return &Cell{stone: stone}, nil
}

func (c *Cell) Put(color Color) error {
	stone, err := NewStone(color)
	if err != nil {
		return fmt.Errorf("failed to put: %w", err)
	}
	c.stone = stone
	return nil
}

func (c *Cell) Draw() {
	var stone string
	if c.stone == nil {
		stone = " "
	} else {
		stone = c.stone.Get()
	}
	fmt.Printf("|%s", stone)
}

func (c *Cell) Reverse() {
	if c.stone != nil {
		c.stone.Reverse()
	}
}

func (c *Cell) IsNone() bool {
	return c.stone == nil
}

func (c *Cell) IsBlack() bool {
	switch {
	case c.stone == nil:
		return false
	case c.stone.color == Black:
		return true
	default:
		return false
	}
}

func (c *Cell) IsWhite() bool {
	switch {
	case c.stone == nil:
		return false
	case c.stone.color == White:
		return true
	default:
		return false
	}
}

type Address struct {
	x int
	y int
}

func NewAddress(x int, y int) *Address {
	return &Address{x: x, y: y}
}

func (a *Address) Valid() error {
	if a.x >= 0 && a.x <= 7 && a.y >= 0 && a.y <= 7 {
		return nil
	} else {
		return fmt.Errorf(`"%v" is invalid as address`, a)
	}
}

func (a *Address) String() string {
	return fmt.Sprintf("x: %d, y: %d", a.x, a.y)
}

func (a *Address) X() int {
	return a.x
}

func (a *Address) Y() int {
	return a.y
}

type Board struct {
	board [][]*Cell
}

func NewBoard() (*Board, error) {
	var board [][]*Cell
	for i := 0; i < 8; i++ {
		var row []*Cell
		for j := 0; j < 8; j++ {
			cell, err := NewCell(nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create board: %w", err)
			}
			row = append(row, cell)
		}
		board = append(board, row)
	}
	if err := board[3][3].Put(Black); err != nil {
		return nil, fmt.Errorf("failed to create board: %w", err)
	}
	if err := board[3][4].Put(White); err != nil {
		return nil, fmt.Errorf("failed to create board: %w", err)
	}
	if err := board[4][3].Put(White); err != nil {
		return nil, fmt.Errorf("failed to create board: %w", err)
	}
	if err := board[4][4].Put(Black); err != nil {
		return nil, fmt.Errorf("failed to create board: %w", err)
	}
	return &Board{board: board}, nil
}

func (b *Board) Draw() {
	fmt.Println("  0 1 2 3 4 5 6 7")
	for i := 0; i < 8; i++ {
		fmt.Print(i)
		for j := 0; j < 8; j++ {
			b.board[i][j].Draw()
		}
		fmt.Println("|")
	}
}

func (b *Board) refCell(address *Address) *Cell {
	return b.board[address.Y()][address.X()]
}

func (b *Board) Put(color Color, address *Address) error {
	targets := b.search(color, address)
	if len(targets) == 0 {
		return fmt.Errorf(`"%s" cannot reverse opponent's stone`, address)
	}
	if err := b.refCell(address).Put(color); err != nil {
		return fmt.Errorf("failed to put stone: %w", err)
	}
	for _, v := range targets {
		b.refCell(v).Reverse()
	}
	return nil
}

func (b *Board) search(color Color, startPoint *Address) []*Address {
	// メモ: クロージャーを再起関数にするためには、関数を定義する前にvarで型宣言する
	// メモ: クロージャーは無名関数で利用可能
	var searchNext func(current *Address, list []*Address, next func(address *Address) *Address) []*Address
	searchNext = func(current *Address, list []*Address, next func(address *Address) *Address) []*Address {
		nextAddress := next(current)
		if err := nextAddress.Valid(); err != nil {
			return nil
		}
		nextCell := b.refCell(nextAddress)
		if nextCell.IsNone() {
			return nil
		}
		if (nextCell.IsBlack() && color == White) || (nextCell.IsWhite() && color == Black) {
			list = append(list, nextAddress)
			return searchNext(nextAddress, list, next)
		}
		return list
	}

	results := []*Address{}
	results = append(results, searchNext(startPoint, []*Address{}, func(a *Address) *Address { return NewAddress(a.X(), a.Y()-1) })...)
	results = append(results, searchNext(startPoint, []*Address{}, func(a *Address) *Address { return NewAddress(a.X(), a.Y()+1) })...)
	results = append(results, searchNext(startPoint, []*Address{}, func(a *Address) *Address { return NewAddress(a.X()-1, a.Y()) })...)
	results = append(results, searchNext(startPoint, []*Address{}, func(a *Address) *Address { return NewAddress(a.X()+1, a.Y()) })...)
	results = append(results, searchNext(startPoint, []*Address{}, func(a *Address) *Address { return NewAddress(a.X()-1, a.Y()-1) })...)
	results = append(results, searchNext(startPoint, []*Address{}, func(a *Address) *Address { return NewAddress(a.X()+1, a.Y()-1) })...)
	results = append(results, searchNext(startPoint, []*Address{}, func(a *Address) *Address { return NewAddress(a.X()-1, a.Y()+1) })...)
	results = append(results, searchNext(startPoint, []*Address{}, func(a *Address) *Address { return NewAddress(a.X()+1, a.Y()+1) })...)

	return results
}

func (b *Board) CountWhite() int {
	count := 0
	for _, row := range b.board {
		for _, cell := range row {
			if cell.IsWhite() {
				count++
			}
		}
	}
	return count
}

func (b *Board) CountBlack() int {
	count := 0
	for _, row := range b.board {
		for _, cell := range row {
			if cell.IsBlack() {
				count++
			}
		}
	}
	return count
}
