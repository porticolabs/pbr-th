package main

import (
	"strings"
	"encoding/json"

    log "github.com/sirupsen/logrus"

	"github.com/dghubble/go-twitter/twitter"
)

// Creates a Demux object to process
//  the tweets from the stream
func createDemux() (twitter.SwitchDemux){
    demux := twitter.NewSwitchDemux()
    demux.Tweet = func(tweet *twitter.Tweet) { processTweet(tweet) }
    return demux
}

// Analyse the tweet data, checking if it should be filtered
//   or published into the Redis queue
func processTweet(tweet *twitter.Tweet) {
	log.WithFields(log.Fields{
		"tweetID":    tweet.ID,
		"tweetUser":    tweet.User.ScreenName,
	}).Debug("Received Tweet")
	// Get the right text
	tweetText := getTweetText(tweet)

	// Check if a Sample Stream is being used, in which case
	//  it won't filter any tweets, otherwise, will run a few
	//  checks to decide if the tweet should continue its way
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

// Run a few checks to decide if the tweet should continue its way
func shouldFilterTweet(tweet *twitter.Tweet) (bool) {
	filterTweet := false

	// Avoid processing own tweets, retweets, 
	//  and quoted tweets without the hashtag
	if ( tweet.User.ScreenName == user.ScreenName ||
		tweet.RetweetedStatus != nil ||
		!strings.Contains(strings.ToLower(getTweetText(tweet)), strings.ToLower(twitterHashtag))) { 
		log.Println(tweet.RetweetedStatus != nil)
		filterTweet = true
	}
	return filterTweet
}

// Tweet that are more than 144 characters long, use a secundary
//   Tweet structure called 'ExtendedTweet' which contains the
//   complete text, and not just the trucated one.
// This functions get returns the tweet complete text wheter if
//   it's extended or not.
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