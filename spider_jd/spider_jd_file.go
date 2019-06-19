package main

import (
  "fmt"
  "log"
  "net/http"
  // "net/url"
  "strings"
  "github.com/PuerkitoBio/goquery"
  "github.com/axgle/mahonia"
  // "regexp"
  // "strconv"
  // "database/sql"
  // _ "github.com/go-sql-driver/mysql"

  "time"
  // "unicode/utf8"
  // "unicode"
  // "io/ioutil"
  // "encoding/json"
  "flag"
  "os"
  "bufio"
  "io"
  "math/rand"
)

// var db = &sql.DB{}
var domain = "https://list.jd.com"
var spiderResultLog = "result.log"
var f *os.File

func ConvertToString(src string, srcCode string, tagCode string) string {
    srcCoder := mahonia.NewDecoder(srcCode)
    srcResult := srcCoder.ConvertString(src)
    tagCoder := mahonia.NewDecoder(tagCode)
    _, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
    result := string(cdata)
    return result
}

func productDetail(url string,img string,cat_id string,brandId string) {
  client := &http.Client{}
  if url[0:1] == "/" {
    url = "https:"+url
  }
  req,err := http.NewRequest("GET",url,nil)

  agent := []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:64.0) Gecko/20100101 Firefox/64.0"}
  agent = append(agent,"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36")
  rand.Seed(time.Now().UnixNano())
  req.Header.Add("user-agent",agent[rand.Intn(2)])
  res, err := client.Do(req)

  if err != nil {
    return
    // log.Fatal(err)
  }

  defer res.Body.Close()
  if res.StatusCode != 200 {
    return
    // log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  }

  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    return
    // log.Fatal(err)
  }

  var model string;
  var model_net string;

  doc.Find(".Ptable-item>h3").Each(func(i int , s *goquery.Selection) {
    h3,_ := s.Html()
    h3    = ConvertToString(h3, "gbk", "utf-8")
    if strings.EqualFold(h3,"主体") || strings.EqualFold(h3,"主体参数")|| strings.EqualFold(h3,"型号信息"){
      s.Parent().Find("dl>dt").Each(func(j int , sel *goquery.Selection) {
        dt,_ := sel.Html()
        utfdt := ConvertToString(dt, "gbk", "utf-8")

        if strings.EqualFold(utfdt,"型号") || strings.EqualFold(utfdt,"产品型号"){
          tips,_ := sel.Parent().Find(".Ptable-tips").Html()
          if tips == "" {
            model,_ = sel.Parent().Find("dd").Eq(0).Html()  
          }else{
            model,_ = sel.Parent().Find("dd").Eq(1).Html()
          }

          return 
        }

        if strings.EqualFold(utfdt,"入网型号") || strings.EqualFold(utfdt,"认证型号") {
          tips,_ := sel.Parent().Find(".Ptable-tips").Html()
          if tips == "" {
            model_net,_ = sel.Parent().Find("dd").Eq(0).Html()  
          }else{
            model_net,_ = sel.Parent().Find("dd").Eq(1).Html()
          }

          return 
        }
      })
      return 
    }
  })
  
  title,_ := doc.Find(".detail .tab-con .parameter2>li").Eq(0).Attr("title")

  // reg := regexp.MustCompile(`[\p{Han}]+`)//查找连续的汉字
  if title == "" {
    return
  }
  
  title = ConvertToString(title, "gbk", "utf-8")
  if model != "" {
    model = ConvertToString(model, "gbk", "utf-8")
  }
  if model_net != "" {
    model_net = ConvertToString(model_net, "gbk", "utf-8")
  }

  line := brandId+"\t"+cat_id+"\t"+title+"\t"+model+"\t"+model_net+"\t"+img
  
  if _, err := f.Write([]byte(line+"\n")); err != nil {
    log.Fatal(err)
  }

  // fmt.Println(line)
}

func detailPageList(url string,brandId string,cat_id string) {
  // Request the HTML page.
  client := &http.Client{}

  req,err := http.NewRequest("GET",url,nil)
  req.Header.Add("user-agent","Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36")
  res, err := client.Do(req)

  if err != nil {
    return
    // log.Fatal(err)
  }

  defer res.Body.Close()
  if res.StatusCode != 200 {
    return
    // log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  }

  // Load the HTML document
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    return
    // log.Fatal(err)
  }
  // var urls []string
  // chkModel := make(map[string]string)
  // Find the review items
  doc.Find("#plist .gl-item").Each(func(i int, s *goquery.Selection) {
    // fmt.Printf("-----------------:%d\n",i)
    // For each item found, get the band and title
    href,_ := s.Find(".p-img a").Attr("href")
    imgSrc,_ := s.Find(".p-img a>img").Attr("src")
    if imgSrc == "" {
      imgSrc,_ = s.Find(".p-img a>img").Attr("data-lazy-img")
    }
    // println(href,imgSrc)
    // urls = append(urls, href)
    // println(href)
    productDetail(href,imgSrc,cat_id,brandId)
    
  })

  nextpage, has := doc.Find("#J_bottomPage .pn-next").Attr("href")
  if !strings.Contains(nextpage,"ev=") {
    return
  }

  if has {
    n_p_url := domain+nextpage
    fmt.Println(n_p_url)
    detailPageList(n_p_url,brandId,cat_id)
  }else{
    return 
  }
  println(nextpage, has)
}


func spider(urlFile string) {
  // file, err := os.Open("spider_jd_url.log")
  file, err := os.Open(urlFile)
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
    // 品牌ID 京东品牌ID 京东分类URL 类别ID
    tmp := strings.Fields(string(line))
    
    jd_url := tmp[2]
    if strings.Contains(jd_url,"ev=") {
      jd_url = jd_url+"%40exbrand_"+tmp[1]
    }else{
      jd_url = jd_url+"&ev=exbrand_"+tmp[1]
    }

    detailPageList(jd_url,tmp[0],tmp[3])
    println(tmp[0],tmp[1],jd_url,tmp[3])
    
    time.Sleep(time.Duration(5) * time.Second)
  }
  println("END")
}

func main() {
  os.Remove(spiderResultLog)
  
  var err error
  f, err = os.OpenFile(spiderResultLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  defer f.Close()
  if err != nil { 
    log.Fatal(err)
  }
  
  flag.Parse()
  var urlFile string
  if len(flag.Args()) == 0 {
    urlFile = "spider_jd_url.log"
  }else{
    urlFile = flag.Args()[0]
  }

  spider(urlFile)

  /*
  flag.Parse()
  item_url := flag.Args()[0]
  _, err := url.Parse(item_url)

  if err != nil {
    panic(err)
  }

  f, err = os.OpenFile(spiderResultLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  defer f.Close()
  if err != nil { 
    log.Fatal(err)
  }

  productDetail(item_url,"img")
  */
}






