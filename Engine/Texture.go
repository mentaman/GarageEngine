package Engine

import (
	"errors"
	"github.com/banthar/gl"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"log"
	//"reflect"
)

var (  
	CustomColorModels = make(map[color.Model]*GLColorModel)
)

type AlignType byte

const (
	AlignLeft = AlignType(0)
	AlignCenter = AlignType(1)
	AlignRight = AlignType(2)
	
	AlignTopLeft = AlignType(4|AlignLeft)
	AlignTopCenter = AlignType(4|AlignCenter)
	AlignTopRight = AlignType(4|AlignRight)
	
	AlignBottomLeft = AlignType(8|AlignLeft)
	AlignBottomCenter = AlignType(8|AlignCenter)
	AlignBottomRight = AlignType(8|AlignRight)
)

func Align(typ AlignType) Vector {
	vect := NewVector2(0,-0.5)
	switch {
		case typ & AlignLeft != 0:
			vect.X = 0
		case typ &  AlignCenter != 0:
			vect.X = -0.5
		case typ &  AlignRight != 0:
			vect.X = -1
	}
	switch {		
		case typ&4 != 0:
			vect.Y = -1
		case typ&8 != 0:
			vect.Y = 0
	}
	return vect
}

type EngineColorModel interface {
	color.Model
	Data() interface{}
}

type GLColorModel struct {
	InternalFormat int
	Type           gl.GLenum
	Format         gl.GLenum
	Target		   gl.GLenum
	PixelBytesSize int
	Model	       EngineColorModel
}

type GLTexture interface {
	GLTexture() gl.Texture
	Height()  int
	Width()   int
	Bind()
}

type Texture struct {
	handle         gl.Texture
	readOnly       bool
	data           interface{}
	format         gl.GLenum
	internalFormat int
	target         gl.GLenum
	width          int
	height         int	
}

func (t *Texture) GLTexture() gl.Texture{
	return t.handle
}

func (t *Texture) Height() int{
	return t.height
}

func (t *Texture) Width() int{
	return t.width
}

func LoadTexture(path string) (tex *Texture, err error) {
	img, e := LoadImage(path)
	if e != nil {
		return nil, e
	}
	return LoadTextureFromImage(img)
}

func LoadImage(path string) (img image.Image, err error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	img, _, e = image.Decode(f)
	if e != nil {
		return nil, e
	}
	return img, nil 
}

func LoadImageQuiet(path string) (img image.Image) {
	f, e := os.Open(path)
	if e != nil {
		log.Println(e)
		return nil
	}
	img, _, e = image.Decode(f)
	if e != nil {
		log.Println(e)
		return nil
	}
	return img
}

func LoadTextureFromImage(image image.Image) (tex *Texture, err error) {
	/*
		val := reflect.ValueOf(image)
		pixs := val.FieldByName("Pix")
		if pixs.IsValid() {
			return
		} else {

		}
		return nil, false, nil
	*/
	
	internalFormat, typ, format, target, e := ColorModelToGLTypes(image.ColorModel())
	if e != nil {
		return nil, nil
	}
	data, e := ImageData(image)
	if e != nil {
		return nil, nil
	} 
	return NewTexture2(data, image.Bounds().Dx(), image.Bounds().Dy(), target, internalFormat, typ, format), nil
}


func ColorModelToGLTypes(model color.Model) (internalFormat int, typ gl.GLenum, format gl.GLenum, target gl.GLenum, err error) {
	switch model {
	case color.RGBAModel, color.NRGBAModel:
		return gl.RGBA8, gl.RGBA, gl.UNSIGNED_BYTE, gl.TEXTURE_2D, nil
	case color.RGBA64Model, color.NRGBAModel:
		return gl.RGBA16, gl.RGBA, gl.UNSIGNED_SHORT, gl.TEXTURE_2D, nil
	case color.AlphaModel:
		return gl.ALPHA, gl.ALPHA, gl.UNSIGNED_BYTE, gl.TEXTURE_2D, nil
	case color.Alpha16Model:
		return gl.ALPHA16, gl.ALPHA, gl.UNSIGNED_SHORT, gl.TEXTURE_2D, nil
	case color.GrayModel:
		return gl.LUMINANCE, gl.LUMINANCE, gl.UNSIGNED_BYTE, gl.TEXTURE_2D, nil
	case color.Gray16Model:
		return gl.LUMINANCE16, gl.LUMINANCE, gl.UNSIGNED_SHORT, gl.TEXTURE_2D, nil
	case color.YCbCrModel:
		return gl.RGB8, gl.RGB, gl.UNSIGNED_BYTE, gl.TEXTURE_2D, nil
	default:
		m, e := CustomColorModels[model]
		if e {
			return m.InternalFormat, m.Type, m.Format, m.Target, nil
		}
		break
	}
	return 0, 0, 0, 0, errors.New("unsupported format")
}

func ImageData(image image.Image) (data interface{},err error) {
	//
	w := image.Bounds().Dx()
	h := image.Bounds().Dy()
	model := image.ColorModel()
	
	switch model {
	case color.YCbCrModel:
		data := make([]byte, 3*h*w)
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				offset := x + y*w
				r, g, b, _ := image.At(x, y).RGBA()
				data[offset] = byte(r / 257)
				data[offset+1] = byte(g / 257)
				data[offset+2] = byte(b / 257)
			}
		}
		return data, nil
	case color.RGBAModel, color.NRGBAModel:
		data2 := make([]byte, 4*h*w)
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				offset := (x + (y*w))*4
				r, g, b, a := image.At(x, y).RGBA()
				data2[offset] = byte(r / 257)
				data2[offset+1] = byte(g  / 257)
				data2[offset+2] = byte(b  / 257)
				data2[offset+3] = byte(a  / 257 )
			} 
		}
		return data2, nil
	case color.RGBA64Model, color.NRGBA64Model:
		data := make([]byte, 4*h*w)
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				offset := x + y*w
				r, g, b, a := image.At(x, y).RGBA()
				data[offset] = byte(r / 257)
				data[offset+1] = byte(g / 257)
				data[offset+2] = byte(b / 257)
				data[offset+3] = byte(a / 257)
			}
		}
		return data, nil
	default:
		m, e := CustomColorModels[model]
		if e {
			return m.Model.Data(),nil
		}
	}
	return nil, errors.New("unsupported format")
}

func NewRGBTexture(rgbData interface{}, width int, height int) *Texture {
	return NewTexture2(rgbData, width, height, gl.TEXTURE_2D, gl.RGB8, gl.RGB, gl.UNSIGNED_BYTE)
}

func NewRGBATexture(rgbaData interface{}, width int, height int) *Texture {
	return NewTexture2(rgbaData, width, height, gl.TEXTURE_2D, gl.RGBA8, gl.RGBA, gl.UNSIGNED_BYTE)
}

func NewTexture(image image.Image, data interface{}) (texture *Texture,err error) {
	iF, typ, f, t, e := ColorModelToGLTypes(image.ColorModel())
	if e != nil {
		return nil, nil
	}
	return NewTexture2(data, image.Bounds().Dx(), image.Bounds().Dy(), t, iF, typ, f), nil
}

func NewTexture2(data interface{}, width int, height int, target gl.GLenum, internalFormat int, typ gl.GLenum, format gl.GLenum) *Texture {
	a := gl.GenTexture()
	a.Bind(target)
	gl.TexImage2D(target, 0, internalFormat, width, height, 0, typ, format, data)
	gl.TexParameteri(target, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(target, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(target, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
    gl.TexParameteri(target, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
    //gl.TexParameteri(target, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	//gl.TexParameteri(target, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	a.Unbind(target)
	
	return &Texture{a,false, data, format, internalFormat, target,width,height}
}

func NewTextureEmpty(width int, height int, model color.Model) *Texture {
	internalFormat, typ, format, target, e := ColorModelToGLTypes(model)
	if e != nil {
		return nil
	}
	a := gl.GenTexture()
	a.Bind(target)
	gl.TexImage2D(target, 0, internalFormat, width, height, 0, typ, format, nil)
	gl.TexParameteri(target, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(target, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(target, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
    gl.TexParameteri(target, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
    //gl.TexParameteri(target, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	//gl.TexParameteri(target, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	a.Unbind(target)
	
	return &Texture{a,false, nil, format, internalFormat, target,width,height}
}

func (t *Texture) Options(filter, clamp int) {
	t.Bind()
	gl.TexParameteri(t.target, gl.TEXTURE_MIN_FILTER, filter)
	gl.TexParameteri(t.target, gl.TEXTURE_MAG_FILTER, filter)
	gl.TexParameteri(t.target, gl.TEXTURE_WRAP_S, clamp);
    gl.TexParameteri(t.target, gl.TEXTURE_WRAP_T, clamp);
	t.Unbind()
}



func (t *Texture) BuildMipmaps() {
	t.Bind()
	//glu.Build2DMipmaps(t.target, t.internalFormat, t.width, t.height, t.format, t.data)
	t.Unbind()
}

func (t *Texture) Image() image.Image {
	if !t.readOnly {
		return nil
	}
	return nil
}

func (t *Texture) ReadTextureFromGPU(buffer []byte) {
	gl.GetTexImage(t.target, 0, t.format, buffer)
}

func (t *Texture) SetReadOnly() {
	if t.readOnly {
		return
	}

	t.data = nil
	t.readOnly = true
}

func (t *Texture) Bind() {
	t.handle.Bind(t.target)
}

func (t *Texture) Unbind() {
	t.handle.Unbind(t.target)
}

func (t *Texture) Render()	{
	t.Bind()
	
	xratio := float32(t.width) / float32(t.height)
	gl.Begin(gl.QUADS) 
	gl.TexCoord2f(0, 1); gl.Vertex3f(-0.5, -0.5, 1) 
	gl.TexCoord2f(1, 1); gl.Vertex3f((xratio)-0.5, -0.5, 1) 
	gl.TexCoord2f(1, 0); gl.Vertex3f((xratio)-0.5, 0.5, 1) 
	gl.TexCoord2f(0, 0); gl.Vertex3f(-0.5, 0.5, 1) 
	gl.End()
	t.Unbind()
}

func (t *Texture) Release() {
	t.data = nil
	gl.DeleteTextures([]gl.Texture{t.handle})
}
