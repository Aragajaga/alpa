package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

type AudioManager struct {
	assets     map[string]*audio.Player
	currentBGM *audio.Player
	game       *Game
}

func (am *AudioManager) Load(key string, path string) {
	fd, err := os.Open(path)
	if err != nil {
		return
	}

	stream, err := vorbis.DecodeWithSampleRate(44100, fd)
	if err != nil {
		return
	}

	am.assets[key], err = audio.NewPlayer(GetAudioContext(), stream)
	if err != nil {
		return
	}
}

func (am *AudioManager) Play(key string) {
	if am.game.volumeMusic > 0 {
		am.assets[key].Pause()
		am.assets[key].Rewind()
		am.assets[key].Play()
	}
}

func (am *AudioManager) PlayBackgroundMusic(key string) {
	/*
		if am.currentBGM != nil {
			am.currentBGM.Pause()
			am.currentBGM.Rewind()
		}

		am.currentBGM = am.game.audioManager.assets[key]
		am.currentBGM.Pause()
		am.currentBGM.Rewind()
		am.currentBGM.Play()
	*/
}

func NewAudioManager(g *Game) *AudioManager {
	am := new(AudioManager)
	am.game = g
	am.assets = make(map[string]*audio.Player)
	return am
}
