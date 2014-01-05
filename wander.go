package main

import (
    f "fmt"
    "github.com/go-gl/gl"
    glfw "github.com/go-gl/glfw3"
    "encoding/binary"
)

var (
    triangle = []float32{
        0.0, -0.5,
        0.5, -0.5,
        -0.5, -0.5}
)

func main(){
    if !glfw.Init(){
        f.Println("Failed to init glfw")
        panic("Cannot initialize glfw library")
    }
    defer glfw.Terminate()

    //glfw.WindowHint(glfw.DepthBits, 16)
    window, err := glfw.CreateWindow(300, 300, "Wander", nil, nil)
    if err != nil{
        panic(err)
    }

    window.SetFramebufferSizeCallback(reshape)
    window.SetKeyCallback(key)
    window.MakeContextCurrent()
    glfw.SwapInterval(1)
    width, height := window.GetFramebufferSize()
    reshape(window, width, height)

    if gl.Init() != 0 {
        panic("Failed to init GL")
    }
    gl.ClearColor(0.0, 0.45, 0.0, 0.0)

    buf := newBuffer()
    defer buf.Delete()

    prog := newProgram()
    defer prog.Delete()
    f.Println(prog.GetInfoLog())
    attrib := prog.GetAttribLocation("position")

    for !window.ShouldClose() {
        attrib.EnableArray()
        attrib.AttribPointer(2, gl.FLOAT, false, 0, nil)
        draw()
        animate()
        attrib.DisableArray()
        window.SwapBuffers()
        glfw.PollEvents()
    }

}

func key(window *glfw.Window, k glfw.Key, s int, action glfw.Action, mods glfw.ModifierKey) {
    if action != glfw.Press {
        return
    }

    switch glfw.Key(k){
        case glfw.KeyEscape:
            window.SetShouldClose(true);
        default:
            return
    }
}

func reshape(window *glfw.Window, width, height int){
    gl.Viewport(0, 0, width, height)
}

func newBuffer() (buf gl.Buffer){
    buf = gl.GenBuffer()
    buf.Bind(gl.ARRAY_BUFFER)
    gl.BufferData(gl.ARRAY_BUFFER, binary.Size(triangle), triangle, gl.STATIC_DRAW)

    return buf
}
func newShader(shaderType string) (shade gl.Shader){
    var shaderString string

    switch shaderType {
        case "vertex":
            shade = gl.CreateShader(gl.VERTEX_SHADER)
            shaderString = `
            #version 150

            in vec2 position;

            void main(){
                gl_Position = vec4(position, 0.0, 1.0);
            }`
            break;
        case "fragment":
            shade = gl.CreateShader(gl.FRAGMENT_SHADER)
            shaderString = `
            #version 150

            out vec4 outColor;
            
            void main(){
                outColor = vec4(1.0, 1.0, 1.0, 1.0);
            }`
            break;
    }
    shade.Source(shaderString)
    shade.Compile()
    return
}
func newProgram()(prog gl.Program){
    prog = gl.CreateProgram()
    prog.AttachShader(newShader("vertex"))
    prog.AttachShader(newShader("fragment"))
    prog.Link()
    return
}
func draw(){
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.DrawArrays(gl.TRIANGLES, 0, 3)
}

func animate(){
}
