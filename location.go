package main

import (
	"fmt"
	"math"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// bounds of dolnyslask area
var (
	minLon = 14.2130
	minLat = 49.9466
	maxLon = 18.2863
	maxLat = 52.0100
)

const (
	TILESIZE        = 256
	ZOOM            = 14
	LOADDISTANCE    = 3
	PRELOADDISTANCE = 4
	CLEANUPDISTANCE = 10
	CLEANUPTIME     = time.Second * 2
)

type Tile struct {
	Path    string
	Texture rl.Texture2D
	x, y    int
	Loaded  bool
	Loading bool
	mu      sync.RWMutex
}

type TileManager struct {
	tileQueue    chan *Tile
	tiles        map[string]*Tile
	visibleTiles []*Tile
	lastCleanup  time.Time
	mu           sync.RWMutex
}

func NewTileManager() *TileManager {
	return &TileManager{
		tileQueue:    make(chan *Tile, 1024),
		tiles:        make(map[string]*Tile),
		visibleTiles: make([]*Tile, 0, 1024),
		lastCleanup:  time.Now(),
	}
}
func deg2tile(lat, lon float64, zoom int) (int, int) {
	latRad := lat * math.Pi / 180.0
	n := math.Exp2(float64(zoom))
	x := int((lon + 180.0) / 360.0 * n)
	y := int((1.0 - math.Log(math.Tan(latRad)+1.0/math.Cos(latRad))/math.Pi) / 2.0 * n)
	return x, y
}

func validTile(x, y int) bool {
	topLeftX, topLeftY := deg2tile(maxLat, minLon, ZOOM)
	bottomRightX, bottomRightY := deg2tile(minLat, maxLon, ZOOM)

	if x >= topLeftX && x <= bottomRightX && y >= topLeftY && y <= bottomRightY {
		return true
	}
	return false

}

func (tm *TileManager) requestLoad(t *Tile) {
	t.mu.Lock()
	if t.Loading || t.Loaded {
		t.mu.Unlock()
		return
	}
	t.Loading = true
	t.mu.Unlock()

	go func(tile *Tile) {
		if validTile(tile.x, tile.y) {
			tm.tileQueue <- tile
		}

	}(t)
}

func (t *Tile) isReady() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.Loaded
}

func (t *Tile) getTexture() rl.Texture2D {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.Texture
}

func (t *Tile) loadTextureNow() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Texture = rl.LoadTexture(t.Path)
	t.Loaded = true
	t.Loading = false
}

func (t *Tile) unload() {
	t.mu.RLock()
	defer t.mu.RUnlock()
	rl.UnloadTexture(t.Texture)
}

func (tm *TileManager) getTile(x, y int) *Tile {
	key := fmt.Sprintf("file_zoom=14_y=%d_x=%d.png", y, x)
	tm.mu.RLock()
	tile, ok := tm.tiles[key]
	tm.mu.RUnlock()
	if ok {
		return tile
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()
	if tile, ok := tm.tiles[key]; ok {
		return tile
	}
	newTile := &Tile{
		Path: fmt.Sprintf("map_cache/%v", key),
		x:    x,
		y:    y,
	}
	tm.tiles[key] = newTile
	return newTile
}

func (tm *TileManager) setVisibleTiles(cameraX, cameraY float32, screenWidth, screenHeight int) {
	startX := int((cameraX-float32(screenWidth)/2)/TILESIZE) - LOADDISTANCE
	endX := int((cameraX+float32(screenWidth)/2)/TILESIZE) + LOADDISTANCE

	startY := int((cameraY-float32(screenHeight)/2)/TILESIZE) - LOADDISTANCE
	endY := int((cameraY+float32(screenHeight)/2)/TILESIZE) + LOADDISTANCE
	tm.visibleTiles = tm.visibleTiles[:0]

	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			if validTile(x, y) {
				tile := tm.getTile(x, y)
				tm.visibleTiles = append(tm.visibleTiles, tile)
			}
		}

	}

}

func (tm *TileManager) preloadNearbyTiles(cameraX, cameraY float32) {
	centerTileX := int(cameraX / TILESIZE)
	centerTileY := int(cameraY / TILESIZE)

	for y := centerTileY - PRELOADDISTANCE; y <= centerTileY+PRELOADDISTANCE; y++ {
		for x := centerTileX - PRELOADDISTANCE; x <= centerTileX+PRELOADDISTANCE; x++ {
			if validTile(x, y) {
				tile := tm.getTile(x, y)
				if !tile.isReady() {
					tm.requestLoad(tile)
				}
			}
		}
	}
}

func (tm *TileManager) cleanupDistantTiles(cameraX, cameraY float32) {
	if time.Since(tm.lastCleanup) < CLEANUPTIME {
		return
	}
	tm.lastCleanup = time.Now()
	centerTileX := int(cameraX / TILESIZE)
	centerTileY := int(cameraY / TILESIZE)
	tm.mu.Lock()
	defer tm.mu.Unlock()
	for key, tile := range tm.tiles {
		distX := math.Abs(float64(tile.x - centerTileX))
		distY := math.Abs(float64(tile.y - centerTileY))
		if distX > CLEANUPDISTANCE || distY > CLEANUPDISTANCE {
			tile.unload()
			delete(tm.tiles, key)
		}
	}
}

func (tm *TileManager) getLoadedTiles() []*Tile {
	tiles := make([]*Tile, 0, len(tm.visibleTiles))
	for _, tile := range tm.visibleTiles {
		if tile.isReady() {
			tiles = append(tiles, tile)
		}
	}
	return tiles
}

func main() {
	rl.InitWindow(1280, 720, "Title")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	tm := NewTileManager()

	centerLat := (minLat + maxLat) / 2
	centerLon := (minLon + maxLon) / 2

	centerXTile, centerYTile := deg2tile(centerLat, centerLon, ZOOM)
	mapX, mapY := float32(centerXTile*TILESIZE), float32(centerYTile*TILESIZE)

	camera := rl.Camera2D{
		Offset: rl.Vector2{
			X: float32(rl.GetScreenWidth() / 2),
			Y: float32(rl.GetScreenHeight() / 2),
		},
		Target: rl.Vector2{
			X: mapX,
			Y: mapY,
		},
		Rotation: 0,
		Zoom:     1.0,
	}

	tm.preloadNearbyTiles(camera.Target.X, camera.Target.Y)

	var isDragging bool

	for !rl.WindowShouldClose() {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			isDragging = true
		}
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			isDragging = false
		}
		if isDragging {
			delta := rl.GetMouseDelta()
			camera.Target.X -= delta.X
			camera.Target.Y -= delta.Y
		}

		select {
		case tile := <-tm.tileQueue:
			tile.loadTextureNow()
		default:

		}
		tm.cleanupDistantTiles(camera.Target.X, camera.Target.Y)
		tm.setVisibleTiles(camera.Target.X, camera.Target.Y, rl.GetScreenWidth(), rl.GetScreenHeight())
		tm.preloadNearbyTiles(camera.Target.X, camera.Target.Y)

		rl.BeginDrawing()

		rl.ClearBackground(rl.White)
		rl.BeginMode2D(camera)

		tiles := tm.getLoadedTiles()
		for _, tile := range tiles {
			if tile.isReady() {
				rl.DrawTexture(tile.getTexture(),
					int32(tile.x*TILESIZE),
					int32(tile.y*TILESIZE),
					rl.White)
			}
		}
		mapX, mapY := latLonToPixel(51.008056510784286, 16.254980596758454, ZOOM)
		rl.DrawCircle(int32(mapX), int32(mapY), 1, rl.Red)
		texture := rl.LoadTexture("osm/output.png")
		rl.DrawTexture(texture, int32(mapX), int32(mapY), rl.White)
		isOnPin(camera, mapX, mapY)

		rl.EndMode2D()
		rl.EndDrawing()
	}
	tm.mu.Lock()
	defer tm.mu.Unlock()
	for _, tile := range tm.tiles {
		tile.unload()
	}

}

func isOnPin(camera rl.Camera2D, posX, posY float32) bool {
	box := rl.NewRectangle(
		posX,
		posY,
		64,
		64,
	)
	mouseWorldPos := rl.GetScreenToWorld2D(rl.GetMousePosition(), camera)
	if rl.CheckCollisionPointRec(mouseWorldPos, box) {
		notificatiionBox := rl.NewRectangle(
			posX,
			posY-64,
			200, 64)
		rl.DrawRectangle(
			int32(notificatiionBox.X),
			int32(notificatiionBox.Y),
			int32(notificatiionBox.Width),
			int32(notificatiionBox.Height),
			rl.White)
	}
	return true
}

func latLonToPixel(lat, lon float64, zoom int) (float32, float32) {
	latRad := lat * math.Pi / 180.0
	n := math.Exp2(float64(zoom))

	xtile := (lon + 180.0) / 360.0 * n
	ytile := (1.0 - math.Log(math.Tan(latRad)+1.0/math.Cos(latRad))/math.Pi) / 2.0 * n

	mapX := float32(xtile * 256)
	mapY := float32(ytile * 256)
	return mapX, mapY
}
