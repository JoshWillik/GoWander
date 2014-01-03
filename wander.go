package main

import (
    f "fmt"
    "github.com/go-gl/gl"
    glfw "github.com/go-gl/glfw3"
)

func main(){
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
    Init()

    for !window.ShouldClose() {
        draw()
        animate()
        window.SwapBuffers()
        glfw.PollEvents()
    }

}

func Init() {

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

func draw() {
}

func animate() {
}
