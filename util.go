package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func mu(a ...interface{}) []interface{} {
	return a
}

func ColorM_Colorize(cm *ebiten.ColorM, clr color.Color) {
	r, g, b, a := clr.RGBA()
	cm.Scale(float64(r)/255,
		float64(g)/255,
		float64(b)/255,
		float64(a)/255)
}

type Vec2f struct {
	X, Y float64
}

func (a Vec2f) ScaleVec2f(b Vec2f) Vec2f {
	a.X *= b.X
	a.Y *= b.Y
	return a
}

func (a Vec2f) Scale(s float64) Vec2f {
	a.X *= s
	a.Y *= s
	return a
}

func (a Vec2f) Add(b Vec2f) Vec2f {
	a.X += b.X
	a.Y += b.Y
	return a
}

func (a Vec2f) Subtract(b Vec2f) Vec2f {
	a.X -= b.X
	a.Y -= b.Y
	return a
}

func (a Vec2f) Distance() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y)
}

func (a Vec2f) Normalize() Vec2f {
	dist := a.Distance()

	a.X /= dist
	a.Y /= dist
	return a
}

func (a Vec2f) Translate(b Vec2f) Vec2f {
	a.X += b.X
	a.Y += b.Y
	return a
}

func (a Vec2f) Negate() Vec2f {
	a.X = -a.X
	a.Y = -a.Y
	return a
}

func (a Vec2f) ToVec2i() Vec2i {
	return Vec2i{int(a.X), int(a.Y)}
}

type Vec2i struct {
	X, Y int
}

func (a Vec2i) ScaleVec2i(b Vec2i) Vec2i {
	a.X *= b.X
	a.Y *= b.Y
	return a
}

func (a Vec2i) Scale(s float64) Vec2i {
	a.X = int(float64(a.X) * s)
	a.Y = int(float64(a.Y) * s)
	return a
}

func (a Vec2i) Add(b Vec2i) Vec2i {
	a.X += b.X
	a.Y += b.Y
	return a
}

func (a Vec2i) Subtract(b Vec2i) Vec2i {
	a.X -= b.X
	a.Y -= b.Y
	return a
}

func (a Vec2i) Distance() float64 {
	return math.Sqrt(float64(a.X*a.X + a.Y*a.Y))
}

func (a Vec2i) Normalize() Vec2i {
	dist := a.Distance()

	a.X = int(float64(a.X) / dist)
	a.Y = int(float64(a.Y) / dist)
	return a
}

func (a Vec2i) Translate(b Vec2i) Vec2i {
	a.X += b.X
	a.Y += b.Y
	return a
}

func (a Vec2i) Negate() Vec2i {
	a.X = -a.X
	a.Y = -a.Y
	return a
}

func (a Vec2i) ToVec2f() Vec2f {
	return Vec2f{float64(a.X), float64(a.Y)}
}

type Rectf struct {
	p1 Vec2f
	p2 Vec2f
}

func (rc Rectf) ToRecti() Recti {
	return Recti{rc.p1.ToVec2i(), rc.p2.ToVec2i()}
}

func (rc Rectf) Size() (float64, float64) {
	return rc.p2.X - rc.p1.X, rc.p2.Y - rc.p1.Y
}

type Recti struct {
	p1 Vec2i
	p2 Vec2i
}

func (rc Recti) ToRectf() Rectf {
	return Rectf{rc.p1.ToVec2f(), rc.p2.ToVec2f()}
}

func (rc Recti) Size() (int, int) {
	return rc.p2.X - rc.p1.X, rc.p2.Y - rc.p1.Y
}

type HSL struct {
	H, S, L float64
}

func hueToRGB(v1, v2, h float64) float64 {
	if h < 0 {
		h += 1
	}
	if h > 1 {
		h -= 1
	}
	switch {
	case 6*h < 1:
		return (v1 + (v2-v1)*6*h)
	case 2*h < 1:
		return v2
	case 3*h < 2:
		return v1 + (v2-v1)*((2.0/3.0)-h)*6
	}
	return v1
}

func (c HSL) ToRGB() color.RGBA {
	h := c.H
	s := c.S
	l := c.L

	if s == 0 {
		// it's gray
		return color.RGBA{uint8(l * 255), uint8(l * 255), uint8(l * 255), 255}
	}

	var v1, v2 float64
	if l < 0.5 {
		v2 = l * (1 + s)
	} else {
		v2 = (l + s) - (s * l)
	}

	v1 = 2*l - v2

	r := hueToRGB(v1, v2, h+(1.0/3.0))
	g := hueToRGB(v1, v2, h)
	b := hueToRGB(v1, v2, h-(1.0/3.0))

	return color.RGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 255}
}

func SinFade(f float64) float64 {
	return 0.5 + math.Sin(f*(math.Pi/2.0))*0.5
}

func strlen(str string) int {
	strRunes := []rune(str)
	return len(strRunes)
}

func substr(str string, start int, len int) string {
	strRunes := []rune(str)
	return string(strRunes[start : start+len])
}
