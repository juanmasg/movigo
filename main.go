package main

import (
    "fmt"
    "log"
//    "sort"
    //"flag"
    "io"
    "os"
    "readers"
    //"net/url"
    //"github.com/alexflint/go-arg"
)

func main(){

    var err error
    var GetReader func(string) io.Reader

    opts := parseCommandLine()
    log.Printf("%+v", opts)
    log.Println(opts.readfrom.Scheme, "H", opts.readfrom.Host, "P", opts.readfrom.Port)

    if opts.verbose{
        log.SetFlags(log.LstdFlags | log.Lshortfile)
    }

    //packages := map[string]string{
    //    "UTX32": "TDT",
    //    "UTX64": "Extra",
    //}

    area := Area(opts.area)

    //areadfrom
    fromprefix := opts.readfrom.Raw

    if opts.readfrom.Scheme == "udp"{
        fromprefix = ""
        GetReader = readers.GetMulticastReader
    }else if opts.readfrom.Scheme == "udpxy"{
        fromprefix = fmt.Sprintf("http://%s/udp/", opts.readfrom.Host)
        GetReader = readers.GetHttpReader
    }else if opts.readfrom.Scheme != ""{
        log.Fatal("Unknown scheme", opts.readfrom.Raw)
    }else if opts.readfrom.Raw == ""{
        log.Fatal("No input specified")
    }else{
        fromprefix += "/"
        GetReader = readers.GetFilesystemReader
    }

    //savem3u
    var m3uwriter io.Writer

    if opts.savem3u.Raw == "stdout"{
        m3uwriter = os.Stdout
    }else if opts.savem3u.Raw != ""{ 
        m3uwriter, err = os.OpenFile(opts.savem3u.Raw, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0777); if err != nil{
            log.Fatal(err)
        }

        defer m3uwriter.(*os.File).Close()
    }

    //savexmltv
    var xmltvwriter io.Writer

    if opts.savexmltv.Raw == "stdout"{
        xmltvwriter = os.Stdout
    }else if opts.savexmltv.Raw != ""{ 
        xmltvwriter, err = os.OpenFile(opts.savexmltv.Raw, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0777); if err != nil{
            log.Fatal(err)
        }

        defer xmltvwriter.(*os.File).Close()
    }

    //streamacceess
    streamprefix := opts.streamaccess.Raw
    if opts.streamaccess.Scheme == "udpxy"{
        streamprefix = fmt.Sprintf("http://%s/udp/", opts.streamaccess.Host)
    }
    //else keep untouched

    movi := NewMovi(area)
    ok := movi.Scan(GetReader, fromprefix); if !ok{
        log.Fatal("Something went wrong scanning %s", area)
    }

    if opts.listpackages{
        movi.ListPackages()
        return
    }

    if opts.savem3u.Raw != ""{
        //groups := movi.GetChannelGroups(packages)

        //var keys []int
        //for k := range groups{
        //    keys = append(keys, k)
        //}
        //sort.Ints(keys)

        //channels := movi.GetChannelList(nil) //packages)
        //DumpIPTVSimple(channels, "172.16.10.9", 9998)
        //data := DumpGroupsAsIPTVSimple(groups, streamprefix)
        channels := movi.GetUniqueChannels()
        data := DumpIPTVSimple(channels, streamprefix)
        m3uwriter.Write(data)
        log.Printf("Channels written to %+v %s", m3uwriter, opts.savem3u)
    }

    if opts.savexmltv.Raw != ""{
        channels := movi.GetUniqueChannels()
        data := dumpXMLTVEPG(channels)
        xmltvwriter.Write(data)
        log.Printf("XMLTV written to %+v %s", xmltvwriter, opts.savexmltv)
        //groups := movi.GetChannelGroups(nil)
        //var keys []int
        //for k := range groups{
        //    keys = append(keys, k)
        //}
        //sort.Ints(keys)

        //data := DumpXMLTV(groups, 
        //m3uwriter.Write(data)
        //log.Printf("Channels written to %+v %s", writer, opts.savem3u)
    }
}

