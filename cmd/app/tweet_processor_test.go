package main
import (
    "testing"
	"github.com/dghubble/go-twitter/twitter"
)

func init(){
	user = &twitter.User{
		ScreenName: "NeoCba",
	}
}

func TestGetTweetTextShort(t *testing.T){
	testTweet := twitter.Tweet{ 
		CreatedAt: "",
		ID: 123,
		Text: "Hola @porticocba, #quieroscifi",
		Truncated: false,
		User: &twitter.User{ ScreenName: "PepePapa"}, 
    } 
	tweetText := getTweetText(&testTweet)

	if tweetText != testTweet.Text {
		t.Fatalf("Expected text:\n\"%s\"\nbut got:\n\"%s\"", testTweet.Text,tweetText)
	}
}

func TestGetTweetTextLong(t *testing.T){
	testTweet := twitter.Tweet{ 
		CreatedAt: "",
		ID: 124,
		Text: "",
		Truncated: true,
		ExtendedTweet: &twitter.ExtendedTweet{FullText: "Hola @poticocba, necesito mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho scifi #quieroscifi"},
        User: &twitter.User{ ScreenName: "PepePapa"}, 
    } 
	
	tweetText := getTweetText(&testTweet)

	if tweetText != testTweet.ExtendedTweet.FullText {
		t.Fatalf("Expected text:\n\"%s\"\nbut got:\n\"%s\"", testTweet.ExtendedTweet.FullText,tweetText)
	}
}

func TestShouldFilterTweetOK(t *testing.T){
	testTweet := twitter.Tweet{ 
		CreatedAt: "",
		ID: 125,
		Text: "",
		Truncated: true,
		ExtendedTweet: &twitter.ExtendedTweet{FullText: "Hola @poticocba, necesito mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho scifi #quieroscifi"},
        User: &twitter.User{ ScreenName: "PepePapa"}, 
    } 
	if shouldFilterTweet(&testTweet) == true {
		t.Fatal("Tweet was filtered but it shouldn't have.")
	}
}

func TestShouldFilterTweetSameUser(t *testing.T){
	testTweet := twitter.Tweet{ 
		CreatedAt: "",
		ID: 125,
		Text: "",
		Truncated: true,
		ExtendedTweet: &twitter.ExtendedTweet{FullText: "Hola @poticocba, necesito mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho mucho scifi #quieroscifi"},
        User: &twitter.User{ ScreenName: "NeoCba"}, 
    } 
	if shouldFilterTweet(&testTweet) == false {
		t.Fatal("Tweet wasn't filtered but it should have.")
	}
}