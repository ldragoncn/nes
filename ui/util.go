package ui

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/fogleman/nes/nes"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func readKey(window *glfw.Window, key glfw.Key) bool {
	return window.GetKey(key) == glfw.Press
}

func readKeys(window *glfw.Window, turbo bool) [8]bool {
	var result [8]bool
	result[nes.ButtonA] = readKey(window, glfw.KeyZ) || (turbo && readKey(window, glfw.KeyA))
	result[nes.ButtonB] = readKey(window, glfw.KeyX) || (turbo && readKey(window, glfw.KeyS))
	result[nes.ButtonSelect] = readKey(window, glfw.KeyRightShift)
	result[nes.ButtonStart] = readKey(window, glfw.KeyEnter)
	result[nes.ButtonUp] = readKey(window, glfw.KeyUp)
	result[nes.ButtonDown] = readKey(window, glfw.KeyDown)
	result[nes.ButtonLeft] = readKey(window, glfw.KeyLeft)
	result[nes.ButtonRight] = readKey(window, glfw.KeyRight)
	return result
}

func readJoystick(joy glfw.Joystick, turbo bool) [8]bool {
	var result [8]bool
	if !glfw.JoystickPresent(joy) {
		return result
	}
	axes := glfw.GetJoystickAxes(joy)
	buttons := glfw.GetJoystickButtons(joy)
	result[nes.ButtonA] = buttons[0] == 1 || (turbo && buttons[2] == 1)
	result[nes.ButtonB] = buttons[1] == 1 || (turbo && buttons[3] == 1)
	result[nes.ButtonSelect] = buttons[6] == 1
	result[nes.ButtonStart] = buttons[7] == 1
	result[nes.ButtonUp] = axes[1] < -0.5
	result[nes.ButtonDown] = axes[1] > 0.5
	result[nes.ButtonLeft] = axes[0] < -0.5
	result[nes.ButtonRight] = axes[0] > 0.5
	return result
}

func joystickReset(joy glfw.Joystick) bool {
	if !glfw.JoystickPresent(joy) {
		return false
	}
	buttons := glfw.GetJoystickButtons(joy)
	return buttons[4] == 1 && buttons[5] == 1
}

func combineButtons(a, b [8]bool) [8]bool {
	var result [8]bool
	for i := 0; i < 8; i++ {
		result[i] = a[i] || b[i]
	}
	return result
}

func hashFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5.Sum(data)), nil
}

func createTexture() uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return texture
}

func setTexture(im *image.RGBA) {
	size := im.Rect.Size()
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA, int32(size.X), int32(size.Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(im.Pix))
}

func copyImage(src image.Image) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Rect, src, image.ZP, draw.Src)
	return dst
}

func loadPNG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return png.Decode(file)
}