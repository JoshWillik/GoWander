package main

import (
    f "fmt"
    "github.com/go-gl/gl"
    glfw "github.com/go-gl/glfw3"
    "math"
    "time"
)

var (
    seconds = time.Now()
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

    prog := setupProgram()
    defer prog.Delete()
    prog.Use()

    arr := getVertexArray()
    defer arr.Delete()
    arr.Bind()

    setup()
    for !window.ShouldClose() {
        draw()
        animate()
        window.SwapBuffers()
        glfw.PollEvents()
    }

}
func setup(){
    //gl.PointSize(40.0)
}
func setupProgram()(prog gl.Program){
    vertexSource := `
        #version 430 core
        
        layout (location = 0) in vec4 offset;

        const vec4 vertecies[3] = vec4[3](
            vec4(0.25, 0.5, 0.5, 1.0),
            vec4(-0.25, 0.5, 0.5, 1.0),
            vec4(-0.25, -0.5, 0.5, 1.0)
        );
        
        void main(){
            gl_Position = vertecies[gl_VertexID];
        }`
    fragmentSource := `
        #version 430 core
        
        out vec4 color;

        void main(){
            color = vec4(1.0, 0.0, 0.0, 0.0); // red, blue, green, ??
        }`
    vert, frag := gl.CreateShader(gl.VERTEX_SHADER), gl.CreateShader(gl.FRAGMENT_SHADER)
    defer vert.Delete()
    defer frag.Delete()
    vert.Source(vertexSource)
    frag.Source(fragmentSource)
    vert.Compile()
    frag.Compile()

    prog = gl.CreateProgram()
    prog.AttachShader(vert)
    prog.AttachShader(frag)
    prog.Link()
    prog.Use()
    f.Println(prog.GetInfoLog())

    return
}
func getVertexArray()(arr gl.VertexArray){
    arr = gl.GenVertexArray()
    return
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
func draw(){
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.DrawArrays(gl.TRIANGLES, 0, 3)
}

func animate(){
    now := float64(time.Since(seconds))

//    offset = []float64{
//        math.Sin(now),
//        math.Cos(now),
//        0.0,0.0}

    red := gl.GLclampf(math.Sin(now) * 0.25 + 0.75)
    blue := gl.GLclampf(math.Cos(now) * 0.25 + 0.75)
    green := gl.GLclampf(time.Since(seconds))
    _ = green;

    gl.ClearColor(red, blue, 0.2, 0.0)
}
