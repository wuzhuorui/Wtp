package main

import (
    "image"
    "os"
    "fmt"
    "bufio"
    "image/png"
    "github.com/golang/freetype"
    "io/ioutil"
    "log"
    "image/draw"
    "image/color"
    "math/rand"
    "sort"
)

func GetContext(s string) *freetype.Context {
    fontBytes,err := ioutil.ReadFile(s)

    if err != nil{
        log.Println("ReadFile",s,err)
        return nil
    }

    f,err := freetype.ParseFont(fontBytes)
    if err != nil{
        log.Println("ParseFont",err)
        return nil
    }

    c := freetype.NewContext()
    c.SetFont(f)
    c.SetDPI(72)
    c.SetFontSize(26)
    return c
}

func DrawBounds(context *freetype.Context,s string,x int , y int,size float64) (int,int){
    lx,ly := -10000,10000
    context.SetFontSize(size)
    pt := freetype.Pt(lx,ly+ int(context.PointToFixed(size)>>6))
    end_pt , _ := context.DrawString(s,pt)
    rx,ry := int(end_pt.X>> 6),int(end_pt.Y >> 6)
    return rx - lx,ry - ly
}

func DrawString(context *freetype.Context,s string,x int, y int,size float64){
    context.SetFontSize(size)
    pt := freetype.Pt(x,y+ int(context.PointToFixed(size)>>6))
    context.DrawString(s,pt)
}

func GetRandomTable(size int)[]int{
    randomtable := make([]int,size)
    for i := 0 ; i < len(randomtable) ; i++{
        randomtable[i] = i
    }
    for i := 0 ; i < len(randomtable) ; i++{
        j := rand.Intn(len(randomtable))
        randomtable[i],randomtable[j] = randomtable[j],randomtable[i]
    }
    return randomtable
}

type CenterTable struct{
    h,w int
    centertable []int 
}

func (t CenterTable)Len()int{
    return len(t.centertable)
}

func (t CenterTable)Less(i int,j int)bool{
    v1 := t.centertable[i]
    v2 := t.centertable[j]
    dx1 := v1 / t.w
    dy1 := v1 % t.w
    dx2 := v2 / t.w
    dy2 := v2 % t.w
    return (dx1 - t.h/2)*(dx1 - t.h/2)+(dy1 -t.w/2)*(dy1-t.w/2) < (dx2 - t.h/2)*(dx2 - t.h/2)+(dy2 - t.w/2)*(dy2-t.w/2)
}

func (t CenterTable)Swap(i,j int){
    t.centertable[i],t.centertable[j] = t.centertable[j],t.centertable[i]
}

func GetCenterTable(h int,w int)[]int{
    centertable := make([]int,h*w)
    for i := 0 ; i < len(centertable) ; i++{
        centertable[i] = i
    }
    sort.Sort(CenterTable{h,w,centertable})
    return centertable
}

func main(){
    dis := NewDisMap(1000,1000)
    randomtable := GetRandomTable(1000*1000)
    dis.SaveAsPng()
    background := image.NewRGBA(image.Rect(0,0,1000,1000))
    draw.Draw(background,background.Bounds(),image.White,image.ZP,draw.Src)
    context := GetContext("font.ttf")
    context.SetClip(background.Bounds())
    context.SetDst(background)
    context.SetSrc(image.Black)

    for fontsize := float64(52);fontsize > 26; fontsize = fontsize -1{
        time := 100 / fontsize
        for idx := 0 ; idx < 1000*1000 && time > 0; idx++{
            px,py := randomtable[idx] % 1000, randomtable[idx] / 1000
            dx,dy := DrawBounds(context,"你好",px,py,fontsize)
            if dis.IsOk(px,py,dx,dy) {
                DrawString(context,"你好",px,py,fontsize)
                dis.Recal(background,image.Rect(px,py,px+dx,py+dy),color.RGBA{255,255,255,255})
                time--
            }
       }
    }
    outFile,err := os.Create("p1.png")

    if err != nil{
        fmt.Println(err)
        os.Exit(-1)
    }
    defer outFile.Close()
    buff := bufio.NewWriter(outFile)

    err = png.Encode(buff,background)

    if err != nil{
        fmt.Println(err)
        os.Exit(-1)
    }

    err = buff.Flush()
    if err != nil{
        fmt.Println(err)
        os.Exit(-1)
    }
    fmt.Println("Save to 1.png")
}
