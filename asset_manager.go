package main

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type AssetManager struct {
	assets map[string]*ebiten.Image
}

var assetManagerLock = &sync.Mutex{}
var assetManager *AssetManager

func NewAssetManager() *AssetManager {
	am := new(AssetManager)
	am.assets = make(map[string]*ebiten.Image)
	return am
}

func AssetManager_GetInstance() *AssetManager {
	if assetManager == nil {
		assetManagerLock.Lock()
		assetManager = NewAssetManager()
		assetManagerLock.Unlock()
	}

	return assetManager
}

func (am *AssetManager) Load(key string, path string) {
	am.assets[key] = LoadImage(path)
}

func (am *AssetManager) Get(key string) *ebiten.Image {
	return am.assets[key]
}
