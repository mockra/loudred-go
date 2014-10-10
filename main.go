package main

import (
  "github.com/ChimeraCoder/anaconda"
  "os"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "encoding/json"
  "net/http"
)

type Item struct {
  Title string
  URL string
  Name string
}

type Response struct {
  Data struct {
    Children []struct {
      Data Item
    }
  }
}

type Post struct {
  Url string
  Title string
  Name string
}

func main() {
  anaconda.SetConsumerKey(os.Getenv("LR_CONSUMER_KEY"))
  anaconda.SetConsumerSecret(os.Getenv("LR_CONSUMER_SECRET"))
  api := anaconda.NewTwitterApi(os.Getenv("LR_ACCESS_TOKEN"), os.Getenv("LR_TOKEN_SECRET"))
  session, err := mgo.Dial(os.Getenv("LR_MONGO"))
  if err != nil {
    panic(err)
  }
  defer session.Close()

  session.SetMode(mgo.Strong, true)
  c := session.DB("loudred").C("posts")

  subreddits := []string{"globaloffensive", "games", "programming", "node", "ruby",
                 "rails", "javascript"}

  for _, sub := range subreddits {

    url := "http://www.reddit.com/r/" + sub + ".json?limit=2"

    res, _ := http.Get(url)
    defer res.Body.Close()

    resp := new(Response)
    err = json.NewDecoder(res.Body).Decode(resp)

    items := make([]Item, len(resp.Data.Children))

    for i, child := range resp.Data.Children {
      items[i] = child.Data
    }

    for _, item := range items {
      result := Post{}

      c.Find(bson.M{"name": item.Name}).One(&result)

      if result.Name == "" {
        status := item.Title + " #" + sub  + " - " + item.URL
        api.PostTweet(status, nil)
        c.Insert(&Post{item.URL, item.Title, item.Name})
      }
    }
  }
}
