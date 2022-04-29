package model

//SlackResponse - response from Slack
type SlackResponse struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

//Attachment - part of SlackResponse
type Attachment struct {
	Text     string `json:"text"`
	ImageURL string `json:"image_url"`
}

//GiphyResponse - represents (part of) JSON response sent by Giphy API
type GiphyResponse struct {
	Data Data `json:"data"`
}

//Data - high level attribute
type Data struct {
	Title  string `json:"title"`
	Images Images `json:"images"`
}

//Images - Contains downsized format GIF info
type Images struct {
	Downsized Downsized `json:"downsized"`
}

//Downsized - Giphy URL for downsized format GIF
type Downsized struct {
	URL string `json:"url"`
}
