package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type IResourceManager interface {
	LoadImage(string) *ebiten.Image
	LoadSound(string) *audio.Player
	LoadFontJSON(string) *Font
}

type ResourceManager struct {
	IResourceManager
	resources map[string]IResource
}

func (resMan *ResourceManager) LoadImage(path string) *ebiten.Image {
	var image *ebiten.Image

	resource, has := resMan.resources[path]

	if !has {
		img, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			log.Println("[ResourceManager] LoadImage " + path + " failed. " + err.Error())
		}

		image = img
		resMan.resources[path] = image
		return image
	}

	return resource.(*ebiten.Image)
}

func (resMan *ResourceManager) LoadSound(path string) *audio.Player {
	var player *audio.Player

	resource, has := resMan.resources[path]

	if !has {
		fd, err := os.Open(path)
		if err != nil {
			return nil
		}

		stream, err := vorbis.DecodeWithSampleRate(44100, fd)
		if err != nil {
			return nil
		}

		player, err = audio.NewPlayer(GetAudioContext(), stream)
		if err != nil {
			return nil
		}

		resMan.resources[path] = player
		return player
	}

	return resource.(*audio.Player)
}

func (resMan *ResourceManager) LoadFontJSON(path string) *Font {
	resource, has := resMan.resources[path]

	if !has {
		font, err := LoadFontFromJSON(path)
		if err != nil {
			log.Println("[ResourceManager] Failed to load JSON font resource \"" + path + "\": " + err.Error())
			return nil
		}

		resMan.resources[path] = font
		return font
	}

	return resource.(*Font)
}

func InitResourceManager(resMan *ResourceManager) {
	resMan.resources = make(map[string]IResource)
}

func NewResourceManager() *ResourceManager {
	s := new(ResourceManager)
	InitResourceManager(s)
	return s
}

var g_resourceManager IResourceManager

func ResourceManager_GetInstance() IResourceManager {
	if g_resourceManager == nil {
		g_resourceManager = NewResourceManager()
	}

	return g_resourceManager
}
