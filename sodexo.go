package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "time"
  "encoding/json"
  "flag"
  "net/url"
  "html"
)

const PageHeader = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN"
 "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">

<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="fi">

<head>
  <meta http-equiv="Content-Type" content="text/html;charset=utf-8" />
  <meta http-equiv="Content-Language" content="fi" />
  <title>Sodexon päivän ruokalista</title>
</head>

<body>
  <table border="0">
`

const PageFooter = `
  </table>
</body>
</html>
`

type Wrappings struct {
  Courses     []Message `json:"courses"`
}

type Message struct {
  Title_fi    string `json:"title_fi"`
  Category    string `json:"category"`
  Properties  string `json:"properties"`
}

func PrintResult(lista *Wrappings, flag_html bool) {
  if flag_html {
    fmt.Println(PageHeader)
    for i := range lista.Courses {
      m := lista.Courses[i]
      fmt.Printf("      <tr><td>%s</td><td>%s (%s)</td></tr>\n", html.EscapeString(m.Category), html.EscapeString(m.Title_fi), html.EscapeString(m.Properties))
    }
    fmt.Println(PageFooter)
  } else {
    for i := range lista.Courses {
      m := lista.Courses[i]
      fmt.Printf("%s: %s, %s\n", m.Category, m.Title_fi, m.Properties)
    }
  }
}

func main() {

  flag_html := flag.Bool("html", false, "Tulosta HTML-muodossa")
  use_proxy := flag.Bool("use-proxy", false , "Käytä proxya (http://www-proxy.itella.net:880)")
  flag.Parse()

  today :=fmt.Sprintf("%d/%d/%d", time.Now().Year(), time.Now().Month(), time.Now().Day())
  ruokalista_url := "http://www.sodexo.fi/ruokalistat/output/daily_json/521/" + today + "/fi"

  client := &http.Client{}

  if *use_proxy {
    proxyUrl, _ := url.Parse("http://www-proxy.itella.net:880")
    client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
  }
  res, err := client.Get(ruokalista_url)
  if err != nil {
    log.Fatal("Virhe ruokalistaa haettaessa: ", err)
  }
  
  data, err := ioutil.ReadAll(res.Body)
  res.Body.Close()
  if err != nil {
    log.Fatal(err)
  }

  var lista Wrappings
  err = json.Unmarshal(data, &lista)
  if err != nil {
    log.Fatal("Virhe JSON:ia luettaessa: ", err)
  }
  PrintResult(&lista, *flag_html)

}

