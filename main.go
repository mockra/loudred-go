package main

import (
  "github.com/ChimeraCoder/anaconda"
  "os"
)

func main() {
  anaconda.SetConsumerKey(os.Getenv("LR_CONSUMER_KEY"))
  anaconda.SetConsumerSecret(os.Getenv("LR_CONSUMER_SECRET"))
  api := anaconda.NewTwitterApi(os.Getenv("LR_ACCESS_TOKEN"), os.Getenv("LR_TOKEN_SECRET"))
}
