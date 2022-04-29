package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/abhirockzz/awsome-slack-backend/model"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	slackResponseStaticText string = "*With* :heart: *from awsome funcy*"
	slackReqTimestampHeader        = "x-slack-request-timestamp"
	slackSignatureHeader           = "x-slack-signature"

	slackSigningSecretEnvVar = "SLACK_SIGNING_SECRET"
	giphyAPIKEYEnvVar        = "GIPHY_API_KEY"

	giphyRandomAPIEndpoint = "http://api.giphy.com/v1/gifs/random"
	urlFormat              = "%s?tag=%s&api_key=%s"
)

var (
	signingSecret string
	apiKey        string
)

func init() {
	signingSecret = os.Getenv(slackSigningSecretEnvVar)
	apiKey = os.Getenv(giphyAPIKEYEnvVar)

	if signingSecret == "" || apiKey == "" {
		log.Fatal("Required environment variable(s) missing")
	}
}

func main() {
	lambda.Start(Funcy)
}

//implements slash command backend
func Funcy(r events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	log.Println("awsome function triggered")

	//body is base64 encoded. decode it first
	payload, err := base64.StdEncoding.DecodeString(r.Body)
	if err != nil {
		log.Println("base64 decoding failed")
		return events.LambdaFunctionURLResponse{Body: "Failed to process request. Please contact the admin", StatusCode: http.StatusInternalServerError}, nil
	}

	slackTimestamp := r.Headers[slackReqTimestampHeader]
	slackSignature := r.Headers[slackSignatureHeader]

	slackSigningBaseString := "v0:" + slackTimestamp + ":" + string(payload)

	if !matchSignature(slackSignature, signingSecret, slackSigningBaseString) {
		log.Println("Signature did not match!")
		return events.LambdaFunctionURLResponse{Body: "Function was not invoked by Slack", StatusCode: http.StatusForbidden}, nil
	}

	log.Println("slack request verified successfully!")

	//parse the application/x-www-form-urlencoded data sent by Slack
	vals, err := parse(payload)
	if err != nil {
		log.Println("unable to parse data sent by slack", err)
		return events.LambdaFunctionURLResponse{Body: "Failed to process request", StatusCode: http.StatusBadRequest}, nil
	}

	giphyTag := vals.Get("text")
	log.Println("invoking giphy api for keyword", giphyTag)

	giphyResp, err := http.Get(fmt.Sprintf(urlFormat, giphyRandomAPIEndpoint, giphyTag, apiKey))
	if err != nil {
		log.Println("giphy api did not respond", err)
		return events.LambdaFunctionURLResponse{Body: "Failed to process request", StatusCode: http.StatusFailedDependency}, nil

	}

	resp, err := ioutil.ReadAll(giphyResp.Body)
	if err != nil {
		log.Println("could not read giphy response", err)
		return events.LambdaFunctionURLResponse{Body: "Failed to process request", StatusCode: http.StatusInternalServerError}, nil
	}

	var gr model.GiphyResponse
	json.Unmarshal(resp, &gr)
	title := gr.Data.Title
	url := gr.Data.Images.Downsized.URL

	log.Println("giphy response. image url", url)

	slackResponse := model.SlackResponse{Text: slackResponseStaticText, Attachments: []model.Attachment{{Text: title, ImageURL: url}}}

	responseBody, err := json.Marshal(slackResponse)

	if err != nil {
		log.Println("could not marshal giphy response", err)
		return events.LambdaFunctionURLResponse{Body: "Failed to process request", StatusCode: http.StatusInternalServerError}, nil
	}

	//slack needs the content-type to be set explicitly - https://api.slack.com/slash-commands#responding_immediate_response

	respHeader := map[string]string{"Content-Type": "application/json"}
	response := events.LambdaFunctionURLResponse{Headers: respHeader, Body: string(responseBody), StatusCode: http.StatusOK}

	log.Println("sent response to slack")

	return response, nil
}

func matchSignature(slackSignature, signingSecret, slackSigningBaseString string) bool {

	//calculate SHA256 of the slackSigningBaseString using signingSecret
	mac := hmac.New(sha256.New, []byte(signingSecret))
	mac.Write([]byte(slackSigningBaseString))

	//hex encode the SHA256
	calculatedSignature := "v0=" + hex.EncodeToString(mac.Sum(nil))

	match := hmac.Equal([]byte(slackSignature), []byte(calculatedSignature))
	return match
}

//adapted from from net/http/request.go --> func parsePostForm(r *Request) (vs url.Values, err error)
func parse(b []byte) (url.Values, error) {
	vals, e := url.ParseQuery(string(b))
	if e != nil {
		log.Println("unable to parse", e)
		return nil, e
	}
	return vals, nil
}
