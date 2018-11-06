package main

import (
    "os"
    "fmt"
    "io/ioutil"
    "context"
    "net/http"
    "github.com/dghubble/oauth1"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/dghubble/go-twitter/twitter"
)

var config = oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
var token = oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))
var twitterHttpClient = config.Client(oauth1.NoContext, token)
var twitterClient = twitter.NewClient(twitterHttpClient)

var httpClient = &http.Client{}

func main() {
    lambda.Start(Handler)
}

func Handler(ctx context.Context) (string, error) {
    tweetJoke()
    return "All done!", nil
}

func tweetJoke() {
    request, err := http.NewRequest("GET", "https://icanhazdadjoke.com", nil)
    if err != nil {
        fmt.Printf("Http request error: %s", err)
    } else {
        request.Header.Add("Accept", "text/plain")
        response, err := httpClient.Do(request)

        if err != nil {
            fmt.Printf("Http response error: %s", err)
        } else {
            defer response.Body.Close()
            contents, err := ioutil.ReadAll(response.Body)
            if err != nil {
                fmt.Printf("Get Body error: %s", err)
            } else {
                joke := string(contents)
                fmt.Println(joke)
                _, _, err := twitterClient.Statuses.Update(joke, nil)
                if err != nil {
                    fmt.Printf("Error while tweeting: %s", err)
                }
            }
        }
    }
}
