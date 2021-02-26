package main

import (
    log "github.com/sirupsen/logrus"

    "github.com/dghubble/go-twitter/twitter"
    "github.com/dghubble/oauth1"
)

var (
    client *twitter.Client
    user *twitter.User
)

// Credentials stores all of our access/consumer tokens
// and secret keys needed for authentication against
// the twitter REST API.
type TwitterCredentials struct {
    ConsumerKey       string
    ConsumerSecret    string
    AccessToken       string
    AccessTokenSecret string
}


// getClient is a helper function that will return a twitter client
// that we can subsequently use to send tweets, or to stream new tweets
// this will take in a pointer to a Credential struct which will contain
// everything needed to authenticate and return a pointer to a twitter Client
// or an error
func getTwitterClient(creds *TwitterCredentials) (*twitter.Client, *twitter.User, error) {
    // Pass in your consumer key (API Key) and your Consumer Secret (API Secret)
    config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
    // Pass in your Access Token and your Access Token Secret
    token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

    httpClient := config.Client(oauth1.NoContext, token)
    client := twitter.NewClient(httpClient)

    // Verify Credentials
    verifyParams := &twitter.AccountVerifyParams{
        SkipStatus:   twitter.Bool(true),
        IncludeEmail: twitter.Bool(true),
    }

    // we can retrieve the user and verify if the credentials
    // we have used successfully allow us to log in!
    user, _, err := client.Accounts.VerifyCredentials(verifyParams)
    if err != nil {
        return nil, nil, err
    }

    return client, user, nil
}

func loginToTwitter(creds *TwitterCredentials){
	var err error
	client, user, err = getTwitterClient(creds)
    if err != nil {
        log.Warn("Error getting Twitter Client")
        log.Error(err)
    }
	log.Info("Logged in as User: " + user.Name)
}

func createStream() (*twitter.Stream, error){

    if (twitterSampleStream){
        params := &twitter.StreamSampleParams{
            StallWarnings: twitter.Bool(true),
            Language: []string{twitterLanguage},
        }
        return client.Streams.Sample(params)
    } else {
        params := &twitter.StreamFilterParams{
            Track: []string{twitterHashtag},
            Language: []string{twitterLanguage},
            StallWarnings: twitter.Bool(true),
        }
        log.WithFields(log.Fields{
            "hashtag": twitterHashtag,
            "language": twitterLanguage,
        }).Debug("Stream initiated")
        return client.Streams.Filter(params)
    }
}