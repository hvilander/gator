package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	if feedURL == "" {
		return nil, errors.New("Missing URL in fetch feed")
	}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("https Status code: %d", res.StatusCode)
	}

	feed := RSSFeed{}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error in reading body: %w", err)
	}
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, fmt.Errorf("error in unmarshaing: %w", err)
	}

	/*
		I had done the decoder impl before having boots check over my work
		Boots pushed me in the direction of the readall + unmarshall
			decoder := xml.NewDecoder(res.Body)
			if err := decoder.Decode(&feed); err != nil {
				fmt.Println("error decoding")
				return nil, err
			}
	*/

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, rssItem := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(rssItem.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(rssItem.Description)
	}

	return &feed, nil

}
