package main

import (
	"strings"
	"encoding/json"

    log "github.com/sirupsen/logrus"

	"github.com/dghubble/go-twitter/twitter"
)

func createDemux() (twitter.SwitchDemux){
    demux := twitter.NewSwitchDemux()
    demux.Tweet = func(tweet *twitter.Tweet) { processTweet(tweet) }
    return demux
}

func processTweet(tweet *twitter.Tweet) {
	log.WithFields(log.Fields{
		"tweetID":    tweet.ID,
		"tweetUser":    tweet.User.ScreenName,
	}).Debug("Received Tweet")
	// Get the right text
	tweetText := getTweetText(tweet)

	if (!twitterSampleStream || shouldFilterTweet(tweet)) {
		log.WithFields(log.Fields{
			"tweetText":    tweetText,
			"tweetID":    tweet.ID,
			"tweetUser":    tweet.User.ScreenName,
		}).Debug("Filtered tweet")
		return
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

func shouldFilterTweet(tweet *twitter.Tweet) (bool) {
	filterTweet := false

	// Avoid processing own tweets, retweets, 
	//  and quoted tweets without the hashtag
	if ( tweet.User.ScreenName == user.ScreenName ||
		tweet.RetweetedStatus != nil ||
		!strings.Contains(strings.ToLower(getTweetText(tweet)), strings.ToLower(twitterHashtag))) { 
		filterTweet = true
	}

	return filterTweet
}

func getTweetText (tweet *twitter.Tweet) (string) {
	if (tweet.Truncated) {
		log.WithFields(log.Fields{
			"tweetText":    tweet.ExtendedTweet.FullText,
			"tweetID":    tweet.ID,
			"tweetUser":    tweet.User.ScreenName,
		}).Debug("Text truncated, checking Extended Tweet")
		return tweet.ExtendedTweet.FullText
	} else {
		return tweet.Text
	}
}