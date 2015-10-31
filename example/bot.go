package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/mix3/ape-slack"
	"github.com/naoya/go-pit"
)

type tumblr struct {
	Response struct {
		Posts []struct {
			Photos []struct {
				OriginalSize struct {
					Url string `json:"url"`
				} `json:"original_size"`
			} `json:"photos"`
		} `json:"posts"`
		TotalPosts int `json:"total_posts"`
	} `json:"response"`
}

func main() {
	config, err := pit.Get("ape-slack")
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())

	conn := ape.New(config["slack.token"])

	conn.AddAction("ping", "pong するよ", func(e *ape.Event) error {
		e.Reply("pong")
		return nil
	})

	conn.AddAction("echo", "echo するよ", func(e *ape.Event) error {
		e.Reply(strings.Join(e.Command().Args(), " "))
		return nil
	})

	conn.AddAction("error", "error 返すよ", func(e *ape.Event) error {
		return fmt.Errorf("test error")
	})

	conn.AddAction("panic", "panic 起こすよ", func(e *ape.Event) error {
		panic("test panic")
		return nil
	})

	total_posts := 0
	conn.AddAction("zoi", "http://ganbaruzoi.tumblr.com/ から画像をランダムで返すよ", func(e *ape.Event) error {
		offset := 0
		if 0 < total_posts {
			offset = rand.Intn(total_posts/20+1) * 20
		}
		urls := []string{}
		url := fmt.Sprintf(
			"http://api.tumblr.com/v2/blog/ganbaruzoi.tumblr.com/posts/photo?api_key=%s&offset=%d",
			config["tumblr.api.key"],
			offset,
		)
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		var r tumblr
		err = json.NewDecoder(res.Body).Decode(&r)
		if err != nil {
			return err
		}

		total_posts = r.Response.TotalPosts

		for _, post := range r.Response.Posts {
			for _, photo := range post.Photos {
				urls = append(urls, photo.OriginalSize.Url)
			}
		}

		if len(urls) == 0 {
			return fmt.Errorf("画像が見つかりませんでした")
		}

		e.ReplyWithoutPermalink(urls[rand.Intn(len(urls))])

		return nil
	})

	conn.Loop()
}
