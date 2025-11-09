package renderer

import (
	"github.com/zodimo/go-clay/kimi2/clay"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type SDL2Renderer struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	font     *ttf.Font
}

func NewSDL2Renderer(title string, width, height int32) (*SDL2Renderer, error) {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return nil, err
	}
	err = ttf.Init()
	if err != nil {
		return nil, err
	}

	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}

	font, err := ttf.OpenFont("assets/DejaVuSans.ttf", 16)
	if err != nil {
		return nil, err
	}

	return &SDL2Renderer{
		window:   window,
		renderer: renderer,
		font:     font,
	}, nil
}

func (r *SDL2Renderer) Render(commands []clay.RenderCommand) {
	r.renderer.SetDrawColor(30, 30, 30, 255)
	r.renderer.Clear()

	for _, cmd := range commands {
		switch cmd.CommandType {
		case clay.CLAY_RENDER_COMMAND_TYPE_RECTANGLE:
			rect := &sdl.Rect{
				X: int32(cmd.BoundingBox.X),
				Y: int32(cmd.BoundingBox.Y),
				W: int32(cmd.BoundingBox.Width),
				H: int32(cmd.BoundingBox.Height),
			}
			c := cmd.Data.(clay.RectangleRenderData).Color
			color := ClayToSDLColor(c)
			r.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
			r.renderer.FillRect(rect)

		case clay.CLAY_RENDER_COMMAND_TYPE_TEXT:
			c := cmd.Data.(clay.TextRenderData).Color
			color := ClayToSDLColor(c)
			surface, err := r.font.RenderUTF8Solid(cmd.Data.(clay.TextRenderData).StringContents, color)
			if err != nil {
				continue
			}
			defer surface.Free()

			texture, err := r.renderer.CreateTextureFromSurface(surface)
			if err != nil {
				continue
			}
			defer texture.Destroy()

			rect := &sdl.Rect{
				X: int32(cmd.BoundingBox.X),
				Y: int32(cmd.BoundingBox.Y),
				W: surface.W,
				H: surface.H,
			}
			r.renderer.Copy(texture, nil, rect)
		}
	}

	r.renderer.Present()
}

func (r *SDL2Renderer) Close() {
	r.font.Close()
	r.renderer.Destroy()
	r.window.Destroy()
	ttf.Quit()
	sdl.Quit()
}

// ClayToSDLColor converts Clay's float32 RGBA color (0.0-1.0) to SDL's Color (0-255)
func ClayToSDLColor(c clay.Color) sdl.Color {
	return sdl.Color{
		R: uint8(c.R * 255),
		G: uint8(c.G * 255),
		B: uint8(c.B * 255),
		A: uint8(c.A * 255),
	}
}
