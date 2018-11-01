package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/SlyMarbo/rss"
	"github.com/asdine/storm"
)

// YTSubscription is just the URL of each YouTuber you are subscribed to
type YTSubscription struct {
	ID        int
	SubURL    string `storm:"unique"`
	SubTitle  string
	SubStatus string
}

// YTSubscriptionEntry has info about each of the parsed and downloaded "episodes" from YouTube
type YTSubscriptionEntry struct {
	ID           int
	Subscription int    `storm:"index"`
	URL          string `storm:"unique"`
	Title        string
	Date         time.Time
	DropboxURL   string
	FileSize     int
}

func main() {
	db, err := storm.Open("ytpodders.boltdb", storm.AutoIncrement())
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	ytsub := YTSubscription{
		SubURL:    "https://www.youtube.com/channel/UCh8rjWtGCIAbwPrZb3Te8bQ",
		SubTitle:  "Gearist",
		SubStatus: "enabled",
	}

	err = db.Save(&ytsub)

	ytsub = YTSubscription{
		SubURL:    "https://www.youtube.com/user/durianriders",
		SubTitle:  "durianrider",
		SubStatus: "enabled",
	}

	err = db.Save(&ytsub)

	ytsub = YTSubscription{
		SubURL:    "https://www.youtube.com/channel/UCBB7sYb14uBtk8UqSQYc9-w",
		SubTitle:  "Steve Ramsey",
		SubStatus: "enabled",
	}

	err = db.Save(&ytsub)

	var ytSubscriptions []YTSubscription
	err = db.All(&ytSubscriptions)

	fmt.Println(ytSubscriptions)

	for _, subscription := range ytSubscriptions {

		var ytSubscriptionEntries []YTSubscriptionEntry

		fmt.Println(subscription.ID)
		err = db.Find("Subscription", subscription.ID, &ytSubscriptionEntries)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}

		fmt.Println(subscription)
		fmt.Println(ytSubscriptionEntries)

		var feedURL string

		split := strings.Split(subscription.SubURL, "/")
		feedid := split[len(split)-1]
		// fmt.Println(feedid)

		if strings.Contains(subscription.SubURL, "channel") {
			feedURL = "https://www.youtube.com/feeds/videos.xml?channel_id=" + feedid
		} else {
			feedURL = "https://www.youtube.com/feeds/videos.xml?user=" + feedid
		}
		feed, error := rss.Fetch(feedURL)
		if error != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		// TODO: Make limit configurable and handle situation where limit > range
		feedSlice := feed.Items[:5]
		for _, item := range feedSlice {
			fmt.Println(item.Title)
			if RSSEntryInDB(item.Link, ytSubscriptionEntries) == false {

				fmt.Printf("Adding new RSS Entry:   %s \n", item.Title)

				entry := YTSubscriptionEntry{
					Subscription: subscription.ID,
					URL:          item.Link,
					Title:        item.Title,
					Date:         item.Date,
					DropboxURL:   "dunno_yet",
					FileSize:     123456,
				}

				fmt.Println(entry)

				err = db.Save(&entry)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %v\n", err)
				}

			}
		}

	}

	/*
	   // Add all entries to RSSXML struct which will be used to generate the rss.xml file
	   ytAllSubscriptionEntries := []YTSubscriptionEntry{}
	   //  err = db.Select(&ytAllSubscriptionEntries, "SELECT DISTINCT subscription, url, title, date, dropboxurl, filesize FROM subscription_entries ORDER BY date DESC")
	   if err != nil {
	     fmt.Fprintf(os.Stderr, "error: %v\n", err)
	     os.Exit(1)
	   }
	   for _, ytItem := range ytAllSubscriptionEntries {
	     fmt.Println(ytItem)
	   }
	*/
}

// RSSEntryInDB checks if we have already downloaded this "episode"
func RSSEntryInDB(link string, dbentries []YTSubscriptionEntry) bool {
	for _, b := range dbentries {
		if b.URL == link {
			return true
		}
	}
	return false
}
