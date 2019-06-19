package main

import (
  "fmt"
  "log"
  "./util"
  // "net/http"
  // "net/url"
  "strings"
  // "github.com/PuerkitoBio/goquery"
  "github.com/axgle/mahonia"
  "regexp"
  // "strconv"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"

  "time"
  "unicode/utf8"
  // "unicode"
  // "io/ioutil"
  // "encoding/json"
  "flag"
  "os"
  "bufio"
  "io"
)

var db = &sql.DB{}
var domain = "https://list.jd.com"

func ConvertToString(src string, srcCode string, tagCode string) string {
    srcCoder := mahonia.NewDecoder(srcCode)
    srcResult := srcCoder.ConvertString(src)
    tagCoder := mahonia.NewDecoder(tagCode)
    _, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
    result := string(cdata)
    return result
}

func main() {
  
  db,_ = sql.Open("mysql", "username:passwd@tcp(ip:port)/dbname?charset=utf8")
  db.SetConnMaxLifetime(50*time.Second)
  defer db.Close()
  err := db.Ping()
  if err != nil{
     log.Fatalln(err)
  }
 
  
  flag.Parse()
  var logFile string
  if len(flag.Args()) == 0 {
    log.Fatalln("please input file path")
  }else{
    logFile = flag.Args()[0]
  }
  // fmt.Println(logFile)
  
  file, err := os.Open(logFile)

  if err != nil {
    panic(err)
  }
  defer file.Close()

  reader := bufio.NewReader(file)
  for {
    line, _, err := reader.ReadLine()
    if err == io.EOF {
      break
    }

    tmp := strings.Split(string(line),"\t")

    tmp[0] = strings.Replace(tmp[0], "\uFEFF", "", -1)
    tmp[2] = strings.TrimSpace(tmp[2])
    tmp[3] = strings.TrimSpace(tmp[3])
    var id string
    reg := regexp.MustCompile(`[\p{Han}]+`)
    var querySql string
    if reg.FindString(tmp[4]) == "" && utf8.RuneCountInString(tmp[4]) > 3 {
      tmp[4] = strings.TrimSpace(tmp[4])
      querySql = "select id from device_info where brand_id = "+tmp[0]+" and license_model='"+tmp[4]+"'"
    }else if tmp[3] != "" && utf8.RuneCountInString(tmp[3]) > 1 && !strings.Contains(tmp[3],"官网") {
      querySql = "select id from device_info where brand_id = "+tmp[0]+" and model ='"+tmp[3]+"'"
    }else { 
      querySql = "select id from device_info where brand_id = "+tmp[0]+" and title ='"+tmp[2]+"'"
    }
    
    fmt.Println(querySql)
// line := brandId+"\t"+cat_id+"\t"+title+"\t"+model+"\t"+model_net+"\t"+img
    // fmt.Println(tmp[3])
    err = db.QueryRow(querySql).Scan(&id) //db为sql.DB
    if err == sql.ErrNoRows {
      imgSrc := tmp[5]
      if imgSrc == "" {
        imgSrc = tmp[6]
      }
      if imgSrc == "" {
        continue
      }

      logo  := util.Upload(imgSrc)
      fmt.Printf("imgSrc:%s,logo:%s\n",imgSrc,logo)

      if logo == "" {
        if logo == "" {
          continue
        }
        fmt.Println(tmp)
        log.Fatal("logo empty")
      }

      stm,_ := db.Prepare("INSERT INTO device_info (brand_id,cat_id,title,model,license_model,logo,op_user) VALUES (?,?,?,?,?,?,?)")
      ret,_ := stm.Exec(tmp[0],tmp[1],tmp[2],tmp[3],tmp[4],logo,"sys")
      stm.Close()
      if lastInsertId, err := ret.LastInsertId(); nil == err {
        fmt.Println("LastInsertId:", lastInsertId)
      }
    } else {
      fmt.Println(id)
    }  
  }
}






