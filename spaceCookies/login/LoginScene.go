package login

import (
	"fmt"
	"github.com/vova616/garageEngine/engine"
	"github.com/vova616/garageEngine/engine/components"
	//"github.com/vova616/garageEngine/engine/components/tween"
	"github.com/vova616/garageEngine/spaceCookies/game"
	_ "image/jpeg"
	_ "image/png"
	//"gl"  
	"strconv"
	"time"
	//"strings"
	//"math"
	//"github.com/vova616/chipmunk"
	//"github.com/vova616/chipmunk/vect"
	//"image"
	//"image/color"
	//"encoding/json"
	"math/rand"
	//"os"
	//"fmt"
)

type LoginScene struct {
	*engine.SceneData
}

var (
	LoginSceneGeneral *LoginScene

	backgroundTexture *engine.Texture
	button            *engine.Texture
	ArialFont         *engine.Font
	ArialFont2        *engine.Font
)

func CheckError(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

func LoadTextures() {
	var e error
	ArialFont, e = engine.NewFont("./data/Fonts/arial.ttf", 24)
	CheckError(e)
	ArialFont.Texture.SetReadOnly()

	ArialFont2, e = engine.NewFont("./data/Fonts/arial.ttf", 24)
	CheckError(e)
	ArialFont2.Texture.SetReadOnly()

	backgroundTexture, e = engine.LoadTexture("./data/spaceCookies/background.png")
	CheckError(e)

	button, e = engine.LoadTexture("./data/spaceCookies/button.png")
	CheckError(e)
}

func (s *LoginScene) Load() {
	engine.SetTitle("Space Cookies")
	LoadTextures()

	rand.Seed(time.Now().UnixNano())

	LoginSceneGeneral = s

	s.Camera = engine.NewCamera()

	cam := engine.NewGameObject("Camera")
	cam.AddComponent(s.Camera)

	cam.Transform().SetScalef(1, 1)

	background := engine.NewGameObject("Background")
	background.AddComponent(engine.NewSprite(backgroundTexture))
	background.AddComponent(game.NewBackground(background.Sprite))
	background.Sprite.Render = false
	//background.Transform().SetScalef(float32(backgroung.Height()), float32(backgroung.Height()), 1)
	background.Transform().SetScalef(800, 800)
	background.Transform().SetPositionf(0, 0)

	gui := engine.NewGameObject("GUI")
	gui.Transform().SetParent2(cam)

	mouse := engine.NewGameObject("Mouse")
	mouse.AddComponent(engine.NewMouse())
	mouse.Transform().SetParent2(gui)

	FPSDrawer := engine.NewGameObject("FPS")
	FPSDrawer.Transform().SetParent2(gui)
	txt := FPSDrawer.AddComponent(components.NewUIText(ArialFont2, "")).(*components.UIText)
	fps := FPSDrawer.AddComponent(engine.NewFPS()).(*engine.FPS)
	fps.SetAction(func(fps float64) {
		txt.SetString("FPS: " + strconv.FormatFloat(fps, 'f', 2, 32))
	})
	txt.SetAlign(engine.AlignLeft)

	FPSDrawer.Transform().SetPositionf(20, float32(engine.Height)-20)
	FPSDrawer.Transform().SetScalef(20, 20)

	/*
		tween.Create(&tween.Tween{Target: FPSDrawer, From: []float32{1}, To: []float32{100},
			Algo: tween.Linear, Type: tween.Scale, Time: time.Second * 3, Loop: tween.PingPong})

			tween.Create(&tween.Tween{Target: FPSDrawer, From: []float32{400}, To: []float32{500},
				Algo: tween.Linear, Type: tween.Position, Time: time.Second * 3, Loop: tween.PingPong, Format: "y"})

			tween.Create(&tween.Tween{Target: FPSDrawer, From: []float32{0}, To: []float32{180},
				Algo: tween.Linear, Type: tween.Rotation, Time: time.Second * 6, Loop: tween.PingPong})

		txt.SetAlign(engine.AlignCenter)
	*/
	/*
		{
			FPSDrawer := engine.NewGameObject("FPS")
			FPSDrawer.Transform().SetParent2(gui)
			txt := FPSDrawer.AddComponent(components.NewUIText(ArialFont, "")).(*components.UIText)
			fps := FPSDrawer.AddComponent(engine.NewFPS()).(*engine.FPS)
			fps.SetAction(func(fps float32) {
				txt.SetString("FPS: " + strconv.FormatFloat(float64(fps), 'f', 2, 32))
			})
			txt.SetAlign(engine.AlignLeft)

			FPSDrawer.Transform().SetPositionf(20, float32(engine.Height)-500)
			FPSDrawer.Transform().SetScalef(20, 20)
		}
	*/

	//
	tBox := engine.NewGameObject("TextBox")
	tBox.Transform().SetParent2(gui)

	txt2 := tBox.AddComponent(components.NewUIText(ArialFont2, "Type your name: ")).(*components.UIText)
	txt2.SetFocus(false)
	txt2.SetWritable(false)
	txt2.SetAlign(engine.AlignLeft)

	tBox.Transform().SetPositionf(float32(engine.Width)/2-txt2.Width()*20, float32(engine.Height)/2)
	tBox.Transform().SetScalef(20, 20)
	//
	input := engine.NewGameObject("TextBoxInput")
	input.Transform().SetParent2(gui)
	p := tBox.Transform().Position()
	p.X += txt2.Width() * 20
	input.Transform().SetPosition(p)
	input.Transform().SetScalef(20, 20)

	name := input.AddComponent(components.NewUIText(ArialFont2, "")).(*components.UIText)
	name.SetFocus(true)
	name.SetWritable(true)
	name.SetAlign(engine.AlignLeft)
	//
	/*
		{
			input := engine.NewGameObject("TextBoxInput")
			input.Transform().SetParent2(gui)
			p := tBox.Transform().Position()
			p.X += txt2.Width() * 20
			input.Transform().SetPosition(p)
			input.Transform().SetScalef(20, 20)

			name := input.AddComponent(components.NewUIText(ArialFont2, "")).(*components.UIText)
			name.SetFocus(true)
			name.SetWritable(true)
			name.SetAlign(engine.AlignTopCenter)
		}
		{
			input := engine.NewGameObject("TextBoxInput")
			input.Transform().SetParent2(gui)
			p := tBox.Transform().Position()
			p.X += txt2.Width() * 20
			input.Transform().SetPosition(p)
			input.Transform().SetScalef(20, 20)

			name := input.AddComponent(components.NewUIText(ArialFont2, "")).(*components.UIText)
			name.SetFocus(true)
			name.SetWritable(true)
			name.SetAlign(engine.AlignBottomRight)
		}
	*/
	//
	errLabel := engine.NewGameObject("TextBoxInput")
	errLabel.Transform().SetParent2(gui)
	errLabel.Transform().SetPositionf(float32(engine.Width)/2, float32(engine.Height)/2-100)
	errLabel.Transform().SetScalef(24, 24)

	errLabelTxt := errLabel.AddComponent(components.NewUIText(ArialFont2, "")).(*components.UIText)
	errLabelTxt.SetFocus(false)
	errLabelTxt.SetWritable(false)
	errLabelTxt.SetAlign(engine.AlignCenter)
	errLabelTxt.Color = engine.Vector{1, 1, 1}
	//
	LoginButton := engine.NewGameObject("LoginButton")
	LoginButton.Transform().SetParent2(cam)
	LoginButton.Transform().SetPositionf(float32(engine.Width)/2, float32(engine.Height)/2-50)
	LoginButton.AddComponent(engine.NewSprite(button))
	LoginButton.AddComponent(engine.NewPhysics(false, 1, 1))
	LoginButton.Physics.Shape.IsSensor = true
	LoginButton.Transform().SetScalef(50, 50)
	LoginButton.Sprite.Color = engine.Vector{0.5, 0.5, 0.5}
	/*
		{
			LoginButton := engine.NewGameObject("LoginButton")
			LoginButton.Transform().SetParent2(cam)
			LoginButton.Transform().SetPositionf(float32(engine.Width)/2, float32(engine.Height)/2-50)
			LoginButton.AddComponent(engine.NewSprite(button))
			LoginButton.AddComponent(engine.NewPhysics(false, 1, 1))
			LoginButton.Physics.Shape.IsSensor = true
			LoginButton.Transform().SetScalef(50, 50)
			LoginButton.Sprite.Color = engine.Vector{0.5, 0.5, 0.5}
			LoginButton.Sprite.SetAlign(engine.AlignTopLeft)
			LoginButton.AddComponent(components.NewUIButton(nil, func(enter bool) {
				if enter {
					LoginButton.Sprite.Color = engine.Vector{0.4, 0.4, 0.4}
				} else {
					LoginButton.Sprite.Color = engine.Vector{0.5, 0.5, 0.5}
				}
			}))
		}
		{
			LoginButton := engine.NewGameObject("LoginButton")
			LoginButton.Transform().SetParent2(cam)
			LoginButton.Transform().SetPositionf(float32(engine.Width)/2, float32(engine.Height)/2-50)
			LoginButton.AddComponent(engine.NewSprite(button))
			LoginButton.AddComponent(engine.NewPhysics(false, 1, 1))
			LoginButton.Physics.Shape.IsSensor = true
			LoginButton.Transform().SetScalef(50, 50)
			LoginButton.Sprite.Color = engine.Vector{0.5, 0.5, 0.5}
			LoginButton.Sprite.SetAlign(engine.AlignBottomRight)
			LoginButton.AddComponent(components.NewUIButton(nil, func(enter bool) {
				if enter {
					LoginButton.Sprite.Color = engine.Vector{0.4, 0.4, 0.4}
				} else {
					LoginButton.Sprite.Color = engine.Vector{0.5, 0.5, 0.5}
				}
			}))
		}
	*/
	loginText := engine.NewGameObject("LoginButtonText")
	loginText.Transform().SetParent2(LoginButton)
	loginText.Transform().SetWorldScalef(24, 24)
	loginText.Transform().SetPositionf(0, 0.1)

	var errChan chan error
	LoginButton.AddComponent(components.NewUIButton(func() {
		if errChan == nil && game.MyClient == nil {
			go game.Connect(name.String(), &errChan)
			errLabelTxt.SetString("Connecting...")
		}
	}, func(enter bool) {
		if enter {
			LoginButton.Sprite.Color = engine.Vector{0.4, 0.4, 0.4}
		} else {
			LoginButton.Sprite.Color = engine.Vector{0.5, 0.5, 0.5}
		}
	}))

	txt2 = loginText.AddComponent(components.NewUIText(ArialFont2, "Log in")).(*components.UIText)
	txt2.SetFocus(false)
	txt2.SetWritable(false)
	txt2.SetAlign(engine.AlignCenter)
	txt2.Color = engine.Vector{1, 1, 1}
	//	

	engine.StartCoroutine(func() {
		for {

			if errChan == nil {
				engine.CoYieldSkip()
				continue
			}
			select {
			case loginErr := <-errChan:
				if loginErr != nil {
					errLabelTxt.SetString(loginErr.Error())
					errChan = nil
				}
			default:

			}
			engine.CoYieldSkip()
		}
	})

	//SPACCCEEEEE
	engine.Space.Gravity.Y = 0
	engine.Space.Iterations = 1

	s.AddGameObject(cam)
	s.AddGameObject(background)

	fmt.Println("LoginScene loaded")
}

func (s *LoginScene) New() engine.Scene {
	gs := new(LoginScene)
	gs.SceneData = engine.NewScene("LoginScene")
	return gs
}
