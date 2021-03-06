package engine

import (
	"github.com/vova616/gl"
)

var (
	defaultBuffer      VBO
	defaultIndexBuffer VBO
	defaultVAO         VAO
)

func initDefaultPlane() {
	defaultBuffer = GenBuffer()

	//Triagles
	data := make([]float32, 20)
	data[0] = -0.5
	data[1] = -0.5
	data[2] = 1

	data[3] = 0.5
	data[4] = -0.5
	data[5] = 1

	data[6] = 0.5
	data[7] = 0.5
	data[8] = 1

	data[9] = -0.5
	data[10] = 0.5
	data[11] = 1

	// UV
	data[12] = 0
	data[13] = 1

	data[14] = 1
	data[15] = 1

	data[16] = 1
	data[17] = 0

	data[18] = 0
	data[19] = 0

	//Setup VAO
	defaultVAO = GenVertexArray()
	defaultVAO.Bind()

	defaultBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, data, gl.STATIC_DRAW)

	gl.AttribLocation.EnableArray(0)
	gl.AttribLocation.EnableArray(1)
	gl.AttribLocation.AttribPointer(0, 3, gl.FLOAT, false, 0, uintptr(0))
	gl.AttribLocation.AttribPointer(1, 2, gl.FLOAT, false, 0, uintptr(12*4))
}

func InsideScreen(ratio float32, position Vector, scale Vector) bool {
	cameraPos := GetScene().SceneBase().Camera.Transform().WorldPosition()

	bigScale := scale.X * ratio
	if scale.Y > bigScale {
		bigScale = scale.Y
	}
	bigScale = -bigScale

	x := (position.X - cameraPos.X) + (bigScale / 2)
	y := (position.Y - cameraPos.Y) + (bigScale / 2)
	if x > float32(Width) || x < bigScale {
		return false
	}
	if y > float32(Height) || y < bigScale {
		return false
	}
	return true
}

func DrawSprite(tex *Texture, uv UV, position Vector, scale Vector, rotation float32, aling AlignType, color Color) {
	if !InsideScreen(uv.Ratio, position, scale) {
		return
	}

	internalMaterial.Begin(nil)

	mp := internalMaterial.ProjMatrix
	mv := internalMaterial.ViewMatrix
	mm := internalMaterial.ModelMatrix
	tx := internalMaterial.Texture
	ac := internalMaterial.AddColor
	ti := internalMaterial.Tiling
	of := internalMaterial.Offset

	defaultVAO.Bind()

	v := Align(aling)
	v.X *= uv.Ratio

	camera := GetScene().SceneBase().Camera
	view := camera.InvertedMatrix()
	model := Identity()
	model.Translate(v.X, v.Y, 0)

	model.Scale(scale.X*uv.Ratio, scale.Y, scale.Z)
	model.Rotate(rotation, 0, 0, -1)
	model.Translate(position.X+0.75, position.Y+0.75, position.Z)

	mv.UniformMatrix4fv(false, view)
	mp.UniformMatrix4f(false, (*[16]float32)(camera.Projection))
	mm.UniformMatrix4fv(false, model)
	ac.Uniform4i(int(color.R), int(color.G), int(color.B), int(color.A))
	ti.Uniform2f(uv.U2-uv.U1, uv.V2-uv.V1)
	of.Uniform2f(uv.U1, uv.V1)

	tx.Uniform1i(0)

	tex.Bind()

	gl.DrawArrays(gl.QUADS, 0, 4)

	internalMaterial.End(nil)
}

func DrawSprites(tex *Texture, uvs []UV, positions []Vector, scales []Vector, rotations []float32, alings []AlignType, colors []Color) {

	internalMaterial.Begin(nil)

	mp := internalMaterial.ProjMatrix
	mv := internalMaterial.ViewMatrix
	mm := internalMaterial.ModelMatrix
	tx := internalMaterial.Texture
	ac := internalMaterial.AddColor
	ti := internalMaterial.Tiling
	of := internalMaterial.Offset

	defaultVAO.Bind()

	camera := GetScene().SceneBase().Camera
	view := camera.InvertedMatrix()
	mv.UniformMatrix4fv(false, view)
	mp.UniformMatrix4f(false, (*[16]float32)(camera.Projection))

	tex.Bind()
	tx.Uniform1i(0)

	for i := 0; i < len(uvs); i++ {

		uv, position, scale := uvs[i], positions[i], scales[i]

		if !InsideScreen(uv.Ratio, position, scale) {
			continue
		}

		rotation, aling, color := rotations[i], alings[i], colors[i]

		v := Align(aling)
		v.X *= uv.Ratio

		model := Identity()
		model.Translate(v.X, v.Y, 0)

		model.Scale(scale.X*uv.Ratio, scale.Y, scale.Z)
		model.Rotate(rotation, 0, 0, -1)
		model.Translate(position.X+0.75, position.Y+0.75, position.Z)

		mm.UniformMatrix4fv(false, model)
		ac.Uniform4f(color.R, color.G, color.B, color.A)
		ti.Uniform2f(uv.U2-uv.U1, uv.V2-uv.V1)
		of.Uniform2f(uv.U1, uv.V1)

		gl.DrawArrays(gl.QUADS, 0, 4)
	}

	internalMaterial.End(nil)
}
