package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"jjjosephhh.com/arranger-clone/constants"
)

type Cell struct {
	row int
	col int
}

func (c *Cell) Draw() {
	x := int32(c.col * constants.CELL_SIDE)
	y := int32(c.row * constants.CELL_SIDE)
	rl.DrawRectangle(x, y, constants.CELL_SIDE, 5, rl.Beige)
	rl.DrawRectangle(x, y, 5, constants.CELL_SIDE, rl.Beige)
}

type PlayerDirection int

const (
	None PlayerDirection = iota
	Up
	Right
	Down
	Left
)

type Player struct {
	x                 int
	y                 int
	nextX             int
	nextY             int
	direction         PlayerDirection
	prevDirection     PlayerDirection
	adjacentPositions map[int]bool
}

func (p *Player) Update() *Player {

	newDirection := None

	if p.direction == None && rl.IsKeyPressed(rl.KeyLeft) {
		p.nextX = p.x - constants.CELL_SIDE
		newDirection = Left
	}
	if p.direction == None && rl.IsKeyPressed(rl.KeyRight) {
		p.nextX = p.x + constants.CELL_SIDE
		newDirection = Right
	}
	if p.direction == None && rl.IsKeyPressed(rl.KeyUp) {
		p.nextY = p.y - constants.CELL_SIDE
		newDirection = Up
	}
	if p.direction == None && rl.IsKeyPressed(rl.KeyDown) {
		p.nextY = p.y + constants.CELL_SIDE
		newDirection = Down
	}

	var playerOther *Player

	if newDirection != None {
		p.direction = newDirection
		p.FindCellsAlongMovement(newDirection)

		switch {
		case p.nextX > (constants.GRID_COLS-1)*constants.CELL_SIDE:
			playerOther = &Player{
				x:             -constants.CELL_SIDE,
				y:             p.y,
				nextX:         0,
				nextY:         p.y,
				direction:     p.direction,
				prevDirection: p.prevDirection,
			}
		case p.nextX < 0:
			playerOther = &Player{
				x:             constants.GRID_COLS * constants.CELL_SIDE,
				y:             p.y,
				nextX:         (constants.GRID_COLS - 1) * constants.CELL_SIDE,
				nextY:         p.y,
				direction:     p.direction,
				prevDirection: p.prevDirection,
			}
		case p.nextY > (constants.GRID_ROWS-1)*constants.CELL_SIDE:
			playerOther = &Player{
				x:             p.x,
				y:             -constants.CELL_SIDE,
				nextX:         p.x,
				nextY:         0,
				direction:     p.direction,
				prevDirection: p.prevDirection,
			}
		case p.nextY < 0:
			playerOther = &Player{
				x:             p.x,
				y:             constants.GRID_ROWS * constants.CELL_SIDE,
				nextX:         p.x,
				nextY:         (constants.GRID_ROWS - 1) * constants.CELL_SIDE,
				direction:     p.direction,
				prevDirection: p.prevDirection,
			}
		}
	}

	if p.direction != None {
		switch p.direction {
		case Up:
			p.y -= 10
			if p.y <= p.nextY {
				p.y = p.nextY
				p.direction = None
				p.ResetCellsAlongMovement()
			}
		case Down:
			p.y += 10
			if p.y >= p.nextY {
				p.y = p.nextY
				p.direction = None
				p.ResetCellsAlongMovement()
			}
		case Left:
			p.x -= 10
			if p.x <= p.nextX {
				p.x = p.nextX
				p.prevDirection = p.direction
				p.direction = None
				p.ResetCellsAlongMovement()
			}
		case Right:
			p.x += 10
			if p.x >= p.nextX {
				p.x = p.nextX
				p.prevDirection = p.direction
				p.direction = None
				p.ResetCellsAlongMovement()
			}
		}
	}

	if p.direction == None {
		switch {
		case p.x > (constants.GRID_COLS-1)*constants.CELL_SIDE:
			p.x = 0
		case p.x < 0:
			p.x = (constants.GRID_COLS - 1) * constants.CELL_SIDE
		case p.y > (constants.GRID_ROWS-1)*constants.CELL_SIDE:
			p.y = 0
		case p.y < 0:
			p.y = (constants.GRID_ROWS - 1) * constants.CELL_SIDE
		}
	}

	return playerOther
}

func (p *Player) FindCellsAlongMovement(newDirection PlayerDirection) {
	keyRow := p.y / constants.CELL_SIDE
	keyCol := p.x / constants.CELL_SIDE
	p.ResetCellsAlongMovement()
	switch {
	case newDirection == Left || newDirection == Right:
		for col := 0; col < constants.GRID_COLS; col++ {
			p.adjacentPositions[p.GetAlongMovementIndex(keyRow, col)] = true
		}
	case newDirection == Up || newDirection == Down:
		for row := 0; row < constants.GRID_ROWS; row++ {
			p.adjacentPositions[p.GetAlongMovementIndex(row, keyCol)] = true
		}
	}
}

func (p *Player) ResetCellsAlongMovement() {
	p.adjacentPositions = map[int]bool{}
}

func (p *Player) IsCellAlongMovement(row, col int) bool {
	key := p.GetAlongMovementIndex(row, col)
	return p.adjacentPositions[key]
}

func (p *Player) GetAlongMovementIndex(row, col int) int {
	return 100000*row + col
}

func NewPlayer(x, y int) *Player {
	return &Player{
		x:             x,
		y:             y,
		nextX:         x,
		nextY:         y,
		direction:     None,
		prevDirection: Right,
	}
}

type SpriteSheet struct {
	frame        int
	frameCounter int
	frameSpeed   int
	rows         int
	columns      int
	spriteWidth  int
	spriteHeight int
	texture      *rl.Texture2D
}

func NewSpriteSheet(filePath string, frameSpeed, rows, columns, spriteWidth, spriteHeight int) *SpriteSheet {
	texture := rl.LoadTexture(filePath)
	return &SpriteSheet{
		frame:        0,
		frameCounter: 0,
		frameSpeed:   frameSpeed,
		rows:         rows,
		columns:      columns,
		spriteWidth:  spriteWidth,
		spriteHeight: spriteHeight,
		texture:      &texture,
	}
}

func (ss *SpriteSheet) UnloadTexture() {
	if ss.texture == nil {
		return
	}
	rl.UnloadTexture(*ss.texture)
}

func (ss *SpriteSheet) UpdateFrame() {
	ss.frameCounter++
	if ss.frameCounter >= (constants.FPS / ss.frameSpeed) {
		ss.frameCounter = 0
		ss.frame++
		if ss.frame >= (ss.rows * ss.columns) {
			ss.frame = 0
		}
	}
}

func (spriteSheet *SpriteSheet) Draw(destRec rl.Rectangle, flipHorizontal bool) {
	// Calculate the row and column of the current frame
	column := spriteSheet.frame % spriteSheet.columns
	row := spriteSheet.frame / spriteSheet.columns

	sourceRec := rl.NewRectangle(
		float32(column*spriteSheet.spriteWidth),
		float32(row*spriteSheet.spriteHeight),
		float32(spriteSheet.spriteWidth),
		float32(spriteSheet.spriteHeight),
	)

	origin := rl.NewVector2(0, 0)

	if flipHorizontal {
		sourceRec.Width = -sourceRec.Width
	}

	rl.DrawTexturePro(*spriteSheet.texture, sourceRec, destRec, origin, 0, rl.White)
}

func main() {

	rl.InitWindow(800, 800, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	spriteSheetPlayerIdle := NewSpriteSheet(
		"assets/free-pixel-art-tiny-hero-sprites/2 Owlet_Monster/Owlet_Monster_Idle_4.png",
		constants.FRAME_SPEED, 1, 4, 32, 32,
	)
	defer spriteSheetPlayerIdle.UnloadTexture()

	spriteSheetPlayerRun := NewSpriteSheet(
		"assets/free-pixel-art-tiny-hero-sprites/2 Owlet_Monster/Owlet_Monster_Run_6.png",
		constants.FRAME_SPEED, 1, 6, 32, 32,
	)
	defer spriteSheetPlayerRun.UnloadTexture()

	player := NewPlayer(0, 0)
	var playerOther *Player

	grassAll := rl.LoadTexture("assets/gdaseljori7b1.png")
	defer rl.UnloadTexture(grassAll)
	grassWidth := 1280
	grassHeight := 1280
	grass01SourceRec := rl.Rectangle{
		X:      0,
		Y:      0,
		Width:  float32(grassWidth / 4),
		Height: float32(grassHeight / 4),
	}

	rl.SetTargetFPS(constants.FPS)

	for !rl.WindowShouldClose() {

		if playerTmp := player.Update(); playerTmp != nil {
			fmt.Println("playerTmp was created")
			playerOther = playerTmp
		}
		if playerOther != nil {
			fmt.Println("playerOther:::", playerOther.x, playerOther.nextX)
			playerOther.Update()
		}

		if player.direction == None {
			spriteSheetPlayerIdle.UpdateFrame()
		} else {
			spriteSheetPlayerRun.UpdateFrame()
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.Beige)
		rl.DrawRectangle(5, 5, constants.SCREEN_WIDTH-10, constants.SCREEN_HEIGHT-10, rl.White)

		for gridRow := 0; gridRow < constants.GRID_ROWS; gridRow++ {
			for gridCol := 0; gridCol < constants.GRID_COLS; gridCol++ {
				grass01DestRec := rl.Rectangle{
					X:      float32(gridCol * constants.CELL_SIDE),
					Y:      float32(gridRow * constants.CELL_SIDE),
					Width:  constants.CELL_SIDE,
					Height: constants.CELL_SIDE,
				}
				grass01Origin := rl.NewVector2(0, 0)
				rl.DrawTexturePro(grassAll, grass01SourceRec, grass01DestRec, grass01Origin, 0, rl.White)
				if player.IsCellAlongMovement(gridRow, gridCol) {
					rl.DrawRectangleLines(
						int32(gridCol*constants.CELL_SIDE),
						int32(gridRow*constants.CELL_SIDE),
						constants.CELL_SIDE,
						constants.CELL_SIDE,
						rl.Beige,
					)
				}
			}
		}

		destRec := rl.NewRectangle(
			float32(player.x+constants.PLAYER_SPRITESHEET_OFFSET),
			float32(player.y+constants.PLAYER_SPRITESHEET_OFFSET),
			float32(constants.CELL_SIDE-2*constants.PLAYER_SPRITESHEET_OFFSET),
			float32(constants.CELL_SIDE-2*constants.PLAYER_SPRITESHEET_OFFSET),
		)
		if player.direction == None {
			spriteSheetPlayerIdle.Draw(destRec, player.prevDirection == Left)
		} else {
			spriteSheetPlayerRun.Draw(destRec, player.direction == Left || (player.direction != Right && player.prevDirection == Left))
		}

		if playerOther != nil {
			destRec = rl.NewRectangle(
				float32(playerOther.x+constants.PLAYER_SPRITESHEET_OFFSET),
				float32(playerOther.y+constants.PLAYER_SPRITESHEET_OFFSET),
				float32(constants.CELL_SIDE-2*constants.PLAYER_SPRITESHEET_OFFSET),
				float32(constants.CELL_SIDE-2*constants.PLAYER_SPRITESHEET_OFFSET),
			)
			if playerOther.direction == None {
				fmt.Println("Should be drawing playerOther idle")
				spriteSheetPlayerIdle.Draw(destRec, playerOther.prevDirection == Left)
			} else {
				fmt.Println("Should be drawing playerOther running")
				spriteSheetPlayerRun.Draw(destRec, playerOther.direction == Left || (playerOther.direction != Right && playerOther.prevDirection == Left))
			}
		}

		for gridRow := 0; gridRow < constants.GRID_ROWS; gridRow++ {
			for gridCol := 0; gridCol < constants.GRID_COLS; gridCol++ {
				if player.IsCellAlongMovement(gridRow, gridCol) {
					rl.DrawRectangleLines(
						int32(gridCol*constants.CELL_SIDE),
						int32(gridRow*constants.CELL_SIDE),
						constants.CELL_SIDE,
						constants.CELL_SIDE,
						rl.Beige,
					)
				}
			}
		}

		rl.EndDrawing()

		if playerOther != nil {
			if playerOther.direction == None {
				playerOther = nil
			}
		}
	}
}
