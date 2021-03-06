package epg

import (
    "fmt"
    "log"
    "io"
)

type EPG struct{
    Files   map[uint16]*EPGFile
}

type EPGFile struct{
    File        *File
    Programs    []*Program
}

func newEPG() *EPG{
    epg := &EPG{}
    epg.Files = make(map[uint16]*EPGFile)
    return epg
}

func ReadMulticastEPG(r io.Reader) *EPG{

    epg := newEPG()

    for{
        data := make([]byte, 1500)
        n, _ := r.Read(data)
        fmt.Printf("% 5d ", n)
        chunk := ParseChunk(data[:n])
        log.Println("SKIP!")

        if chunk.End{ break }
    }

    filedata := make([]byte, 0)
    for len(epg.Files) < 500{ // avoid looping forever
        data := make([]byte, 1500)
        n, _ := r.Read(data)
        //fmt.Printf("% 5d ", n)
        fmt.Printf(".")
        chunk := ParseChunk(data[:n])
        filedata = append(filedata, chunk.Data...)
        if chunk.End{
            file := ParseFile(filedata)
            if file != nil && file.Size > 0 && len(file.Data) != 0{
                //fmt.Printf("FILE! %s: ", file.ServiceURL)
                //fmt.Printf("%d", file.ServiceId)
                _, ok := epg.Files[file.ServiceId]; if ok{
                    break
                }
                programs := ParsePrograms(file.Data)
                //fmt.Printf("%d programs\n ", len(programs))
                fmt.Printf("%d/%d ", file.ServiceId, len(programs))
                epg.Files[file.ServiceId] = &EPGFile{file,programs}
            }
            filedata = make([]byte, 0)
        }
    }
    r.(io.Closer).Close()
    return epg
}
