package main

import (
    "image"
    "os"
    "fmt"
    "bufio"
    "image/png"
)

func main(){
    background := image.NewRGBA(image.Rect(0,0,500,500))
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
