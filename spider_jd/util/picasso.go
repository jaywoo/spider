package util
// package main

import (
    // "fmt"
    "net/http"
    "strings"
    "log"
    "io/ioutil"
    "encoding/json"
    "time"
)

type Data struct {
    KB      string `json:"KB"`
    SIZE    string `json:"SIZE"`
    KEY     string `json:"KEY"`
    URL     []string `json:"URL"`
    SAVE_FLOW_KB  string `json:"SAVE_FLOW_KB"`
}

type Content struct {
    STATUS string `json:"STATUS"`
    ERRORMSG string `json:"ERRORMSG"`
    ERRORNO int32 `json:"ERRORNO"`
    DATA map[string]*Data `json:"DATA"`
}

var picassUrl = ""

func Upload(imgUrl string) string {
    client := &http.Client{}
    if imgUrl[0:1] == "/" {
        imgUrl = "http:"+imgUrl
    }

    param := ""
    
    req,err := http.NewRequest("POST",picassUrl,strings.NewReader(param))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add("user-agent","Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36")
    res, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    body, err := ioutil.ReadAll(res.Body)
    res.Body.Close()

    var mapResult map[string]interface{}
    if err := json.Unmarshal([]byte(body), &mapResult); err != nil {
      log.Fatal(err)  
    }

    var picassoUrl string
    tryTimes := 1
    for{
        if tryTimes > 5 {
            break
        }

        var queryResult []Content
        queryUrl := "".(string)

        resp,err := http.Get(queryUrl)
        
        if err != nil {
            println(err)
        }
        queryBody, err := ioutil.ReadAll(resp.Body)
        resp.Body.Close()
        
        if err := json.Unmarshal([]byte(queryBody), &queryResult); err != nil || len(queryResult) == 0  {
            time.Sleep(time.Duration(100*tryTimes)*time.Millisecond)
            tryTimes = tryTimes+1
            continue
        }
        // fmt.Println(queryResult)
        picassoUrl = ""
        break
    }

    return picassoUrl
    
}

// func main() {
//     p_url := Upload("")
//     fmt.Println(p_url)
// }