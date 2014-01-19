package main

import (
    f "fmt"
    //"github.com/go-gl/gl"
    "github.com/JoshWillik/gl"
    glfw "github.com/go-gl/glfw3"
    "math"
    "time"
    "runtime"
)

var (
    numRendered = 0
    lastDraw = time.Now()
    fps = 60
    seconds = time.Now()
    position gl.AttribLocation
    color gl.AttribLocation
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

    position = prog.GetAttribLocation("offset")
    color  = prog.GetAttribLocation("color")

    setup()
    for !window.ShouldClose() {
        if shouldRender(){
            draw()
        }
        animate()
        window.SwapBuffers()
        glfw.PollEvents()
    }

}
func setup(){
    runtime.LockOSThread()
    gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
    colorOffset := [4]float32{
        1.0,
        0.0,
        0.0,
        0.0 }
    color.Attrib4fv(&colorOffset)
    f.Println(gl.GetString(gl.VERSION))
}
func setupProgram()(prog gl.Program){
    vertexSource := `
        #version 430 core
        
        layout (location = 0) in vec4 offset;
        layout (location = 1) in vec4 color;

        const vec4 vertecies[3] = vec4[3](
            vec4(0.25, 0.5, 0.5, 1.0),
            vec4(-0.25, 0.5, 0.5, 1.0),
            vec4(-0.25, -0.5, 0.5, 1.0)
        );

        out COLOR
        {
            vec4 color;
        } out_color;

        void main(){
            gl_Position = vertecies[gl_VertexID] + offset;

            out_color.color = color;
        }`
    fragmentSource := `
        #version 430 core
        
        in COLOR
        {
            vec4 color;
        } in_color;

        out vec4 color;

        void main(){
            color = in_color.color;
        }`
    tesselationControlSource := `
        #version 430 core

        layout (vertices = 3) out;
        
        void main(void){
            if(gl_InvocationID == 0){
                gl_TessLevelInner[0] = 5.0;
                gl_TessLevelOuter[0] = 5.0;
                gl_TessLevelOuter[1] = 5.0;
                gl_TessLevelOuter[2] = 5.0;
            }
            gl_out[gl_InvocationID].gl_Position = gl_in[gl_InvocationID].gl_Position;
        }
    `
    tesselationEvaluationSource := `
        #version 430 core

        layout (triangles) in;

        void main(void){
            gl_Position = 
                (gl_TessCoord.x * gl_in[0].gl_Position) +
                (gl_TessCoord.y * gl_in[1].gl_Position) +
                (gl_TessCoord.z * gl_in[2].gl_Position);
        }`

    vert := gl.CreateShader(gl.VERTEX_SHADER)
    frag := gl.CreateShader(gl.FRAGMENT_SHADER)
    tessControl := gl.CreateShader(gl.TESS_CONTROL_SHADER)
    tessEval := gl.CreateShader(gl.TESS_EVALUATION_SHADER)

    defer vert.Delete()
    defer frag.Delete()
    defer tessControl.Delete()
    defer tessEval.Delete()

    vert.Source(vertexSource)
    frag.Source(fragmentSource)
    tessControl.Source(tesselationControlSource)
    tessEval.Source(tesselationEvaluationSource)

    vert.Compile()
    frag.Compile()
    tessControl.Compile()
    tessEval.Compile()

    f.Println(vert.GetInfoLog())
    f.Println(frag.GetInfoLog())
    f.Println(tessControl.GetInfoLog())
    f.Println(tessEval.GetInfoLog())

    prog = gl.CreateProgram()

    prog.AttachShader(vert)
    prog.AttachShader(frag)
    prog.AttachShader(tessControl)
    prog.AttachShader(tessEval)
    prog.Link()
    prog.Use()
    f.Println(prog.GetInfoLog())

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
func shouldRender() bool{
    if int(time.Since(lastDraw) / time.Millisecond) >= 1000/fps{
        numRendered++
        lastDraw = time.Now()
        return true
    }

    return false;
}

func animate(){
    now := float64(time.Since(seconds)) / float64(time.Second)

    offset := [4]float32{
        float32(math.Sin(now)),
        float32(math.Cos(now)),
        0.0,0.0}
    position.Attrib4fv(&offset)


    red := gl.GLclampf(math.Sin(now) * 0.25 + 0.75)
    blue := gl.GLclampf(math.Cos(now) * 0.25 + 0.75)
    green := gl.GLclampf(time.Since(seconds))
    _, _, _ = green, blue, red;

    gl.ClearColor(0.0,0.0, 0.0, 0.0)
}
