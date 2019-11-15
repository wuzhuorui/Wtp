package main

import(
    "fmt"
    "image"
    "os"
    "bufio"
    "image/png"
    "image/color"
)

type DisMap struct{
    width,height int
    right_dis []int
    down_dis []int
}

func NewDisMap(w int,h int) *DisMap{
    buf1 := make([]int,w*h)
    buf2 := make([]int,w*h)
    for i := 0 ; i < h ; i++{
        for j := 0 ; j < w ; j++{
            buf1[i*w+j] = w - j
            buf2[i*w+j] = w - i
        }
    }
    return &DisMap{w,h,buf1,buf2}
}

func (dis *DisMap) Bounds() image.Rectangle{
    return image.Rect(0,0,dis.width,dis.height)
}

func (dis *DisMap) At(x int,y int) (*int,*int) {
    return &dis.right_dis[y*dis.width+x],&dis.down_dis[y*dis.width+x]
}

func (dis *DisMap) Set(x int,y int){
}

func (dis *DisMap) SaveAsPng(){
    outFile,_:= os.Create("dismap.png")
    pngimage := image.NewRGBA(image.Rect(0,0,dis.width,dis.height))

    for i := 0 ; i < dis.height ; i++{
        for j:= 0; j < dis.width ; j++{
            b1,b2 := dis.At(i,j)
            if *b1 > 0 && *b2 > 0{
                pngimage.Set(i,j,color.RGBA{255,255,255,255})
            }else{
                pngimage.Set(i,j,color.RGBA{0,0,0,255})
            }
        }
    }
    buff := bufio.NewWriter(outFile)
    png.Encode(buff,pngimage)
    buff.Flush()
}

func (dis *DisMap) Recal(src image.Image,bounds image.Rectangle,bg color.Color){
    bounds = bounds.Intersect(src.Bounds())
    bounds = bounds.Intersect(dis.Bounds())
    fmt.Println(bounds)
    if bounds.Dx() == 0 || bounds.Dy() == 0 {
        return
    }

    for i := bounds.Max.X -1 ; i >= bounds.Min.X ; i--{
        for j := bounds.Max.Y - 1 ; j >= 0 ; j--{
            _,b1 := dis.At(i,j)
            c := src.At(i,j)

            cr,cg,cb,ca := c.RGBA()
            br,bg,bb,ba := bg.RGBA()
            if cr == br && cg == bg && cb == bb && ca == ba {
                prej := j + 1
                if prej < dis.Bounds().Max.Y{
                    _,b2 := dis.At(i,prej)
                    *b1 = *b2 + 1
                }
            }else{
                *b1 = 0
            }
        }
    }
    for i := bounds.Max.X -1 ; i >= 0 ; i--{
        for j := bounds.Max.Y - 1 ; j >= bounds.Min.Y ; j--{
            b1,_ := dis.At(i,j)
            c := src.At(i,j)

            cr,cg,cb,ca := c.RGBA()
            br,bg,bb,ba := bg.RGBA()

            if cr == br && cg == bg && cb == bb && ca == ba {
                prei := i + 1
                if prei < dis.Bounds().Max.X{
                    b2,_ := dis.At(prei,j)
                    *b1 = *b2 + 1
                }
            }else{
                *b1 = 0
            }
        }
    }
    dis.SaveAsPng()
}
