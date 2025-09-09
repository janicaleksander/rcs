package infouserstate

import (
	"fmt"
	"math"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type LocationMap struct {
	width, height    float32
	tm               *TileManager
	camera           rl.Camera2D
	isDraggingCamera bool
}

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
	LOADDISTANCE    = 10
	PRELOADDISTANCE = 10
	CLEANUPDISTANCE = 20
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
	return x >= topLeftX && x <= bottomRightX && y >= topLeftY && y <= bottomRightY
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
func latLonToPixel(lat, lon float64, zoom int) (float32, float32) {
	latRad := lat * math.Pi / 180.0
	n := math.Exp2(float64(zoom))

	xtile := (lon + 180.0) / 360.0 * n
	ytile := (1.0 - math.Log(math.Tan(latRad)+1.0/math.Cos(latRad))/math.Pi) / 2.0 * n

	mapX := float32(xtile * TILESIZE)
	mapY := float32(ytile * TILESIZE)
	return mapX, mapY
}
