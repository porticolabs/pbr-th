package main

import (
	"os"
    "os/signal"
    "syscall"
    "strconv"

    log "github.com/sirupsen/logrus"
)

var (
    version             string // version number
    sha1ver             string // sha1 revision used to build the program
    buildTime           string // when the executable was built
    twitterCreds        TwitterCredentials
    twitterHashtag      string
    twitterLanguage     string
    redisCreds          RedisCredentials
    redisQueue          string
    twitterSampleStream bool
)

func init() {
    // do something here to set environment depending on an environment variable
    // or command-line flag
    if os.Getenv("ENVIRONMENT") == "prod" {
        log.SetFormatter(&log.JSONFormatter{})
      } 

    switch os.Getenv("LOG_LEVEL") {
        case "DEBUG":
            log.SetLevel(log.DebugLevel)
            log.Warn("Log level set to DEBUG")
        case "WARN":
            log.SetLevel(log.WarnLevel)
            log.Warn("Log level set to WARN")
        case "ERROR":
            log.SetLevel(log.ErrorLevel)
        case "FATAL":
            log.SetLevel(log.FatalLevel)
    }
    // Getting the Twitter credentials from env
	twitterCreds = TwitterCredentials{
        AccessToken:       os.Getenv("ACCESS_TOKEN"),
        AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
        ConsumerKey:       os.Getenv("CONSUMER_KEY"),
        ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
    }

    twitterHashtag = os.Getenv("TWITTER_HASHTAG")
    twitterLanguage = os.Getenv("TWITTER_LANGUAGE")

    twitterSampleStream, _ = strconv.ParseBool(os.Getenv("TWITTER_SAMPLE"))
    
    // Getting the Twitter credentials from env
	redisCreds = RedisCredentials{
        Host:     os.Getenv("REDIS_HOST"),
        Protocol: os.Getenv("REDIS_PROTOCOL"),
    }
    redisQueue = os.Getenv("REDIS_QUEUE")
  }

func main() {
    // Startup Information
    log.Infof("Initiating PBR (Twiiter Hashtags) v%s", version)
    log.Infof(" * Commit Hash: %s", sha1ver)
    log.Infof(" * Build Date: %s", buildTime)
	
	log.Debug("Signing in to Twitter.")
	loginToTwitter(&twitterCreds)
    
	log.Debug("Creating Twitter Stream")
	stream, _ := createStream()
    
	log.Debug("Creating Twitter Stream Demux")
	demux := createDemux()
	
    log.Debug("Connecting to Redis")
    loginToRedis(redisCreds)

    log.Debug("Opening Redis Queue")
    openQueue(redisQueue)
	
	log.Debug("Initiating Twitter Stream Demux")
	go demux.HandleChan(stream.Messages)

    // Wait for SIGINT and SIGTERM (HIT CTRL-C)
    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
    log.Println(<-ch)
    
    log.Info("Stream reading stoped")
    stream.Stop()
}