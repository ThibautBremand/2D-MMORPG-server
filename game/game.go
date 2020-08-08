package game

import "server/entity"

// game contains some methods and calculations that are purely game dependant.
// They only rely to the game's logic.

// moveCoordinates calculates the new coordinates of a character given a direction.
func MoveCoordinates(character *entity.CharacterView, direction float64) (int, int) {
	x := character.X
	y := character.Y

	switch direction {
	case 0:
		y = y - 1
	case 1:
		x = x - 1
	case 2:
		y = y + 1
	case 3:
		x = x + 1
	}

	return x, y
}

func TPCoordinates(character *entity.CharacterView, gamemap *entity.Gamemap, direction float64) (int, int) {
	x := character.X
	y := character.Y

	switch direction {
	case 0:
		y = gamemap.Height()
	case 1:
		x = gamemap.Width()
	case 2:
		y = gamemap.EdgeMargin()
	case 3:
		x = gamemap.EdgeMargin()
	}

	return x, y
}
