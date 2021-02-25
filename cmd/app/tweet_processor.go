package main

import (
	"strings"
	"encoding/json"

    log "github.com/sirupsen/logrus"

	"github.com/dghubble/go-twitter/twitter"
)

func createDemux(hashtag string, sample bool) (twitter.SwitchDemux){
    demux := twitter.NewSwitchDemux()
    demux.Tweet = func(tweet *twitter.Tweet) {
		log.WithFields(log.Fields{
			"tweetID":    tweet.ID,
			"tweetUser":    tweet.User.ScreenName,
		}).Debug("Received Tweet")
		// Get the right text
		tweetText := ""

		if (tweet.Truncated) {
			tweetText = tweet.ExtendedTweet.FullText
		} else {
			tweetText = tweet.Text
		}

		if (!sample){
			// Avoid processing own tweets, retweets, 
			//  and quoted tweets without the hashtag
			if ( tweet.User.ScreenName == user.ScreenName ||
				tweet.RetweetedStatus != nil ||
				!strings.Contains(strings.ToLower(tweetText), strings.ToLower(hashtag))) { 
				log.WithFields(log.Fields{
					"tweetText":    tweetText,
					"tweetID":    tweet.ID,
					"tweetUser":    tweet.User.ScreenName,
				}).Debug("Filtered tweet")
				return
			}
		}
		
        tweetBytes, err := json.Marshal(tweet)
		if err != nil {
			log.WithFields(log.Fields{
				"tweetText":    tweetText,
				"tweetID":    tweet.ID,
				"tweetUser":    tweet.User.ScreenName,
			}).Warn("Failed to parse Tweet to Json")
			return
		}

		if err := publishBytesToQueue(tweetBytes); err != nil {
			log.WithFields(log.Fields{
				"tweetText":    tweetText,
				"tweetID":    tweet.ID,
				"tweetUser":    tweet.User.ScreenName,
			}).Warn("Failed to publish message")
			return
		}

        // Received tweet info
		log.WithFields(log.Fields{
			"tweetText":    tweetText,
			"tweetID":    tweet.ID,
			"tweetUser":    tweet.User.ScreenName,
		}).Info("Processed tweet")
	}
    return demux
}