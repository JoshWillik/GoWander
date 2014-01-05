package main

import (
    f "fmt"
    "github.com/go-gl/gl"
    glfw "github.com/go-gl/glfw3"
)

var (
    triangle []float32
    array = []int32{
        1,2,3,4,5,6,7,8,
        9,10,11,12,13,14,15,16}
    slice = array[:]
)

func main(){
    f.Println("%v", cap(array))
    if !glfw.Init(){
        f.Println("Failed to init glfw")
        panic("Cannot initialize glfw library")
    }
    defer glfw.Terminate()

    glfw.WindowHint(glfw.DepthBits, 16)
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

    buf := newBuffer()
    defer buf.Delete()

    for !window.ShouldClose() {
        draw()
        animate()
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

func newBuffer() gl.Buffer{
    triangle = []float32 {
        0.0, 0.5,
        0.5, -0.5,
        -0.5, -0.5 }
    buf := gl.GenBuffer()
    buf.Bind(gl.ARRAY_BUFFER)
    gl.BufferData(gl.ARRAY_BUFFER, len(triangle) * 4, triangle, gl.STATIC_DRAW)
    return buf
}
func draw(){
}

func animate(){
}
