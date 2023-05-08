package main

import "encoding/xml"

type Items []struct {
	Text        string `xml:",chardata"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Content     struct {
		Text        string `xml:",chardata"`
		URL         string `xml:"url,attr"`
		Type        string `xml:"type,attr"`
		Expression  string `xml:"expression,attr"`
		Width       string `xml:"width,attr"`
		Height      string `xml:"height,attr"`
		Description struct {
			Text string `xml:",chardata"`
			Type string `xml:"type,attr"`
		} `xml:"description"`
		Credit struct {
			Text   string `xml:",chardata"`
			Role   string `xml:"role,attr"`
			Scheme string `xml:"scheme,attr"`
		} `xml:"credit"`
	} `xml:"content"`
	Encoded string `xml:"encoded"`
	Guid    string `xml:"guid"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
}

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Content string   `xml:"content,attr"`
	Atom    string   `xml:"atom,attr"`
	Media   string   `xml:"media,attr"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text string `xml:",chardata"`
		Link struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Language    string `xml:"language"`
		Category    string `xml:"category"`
		Copyright   string `xml:"copyright"`
		Image       struct {
			Text  string `xml:",chardata"`
			Title string `xml:"title"`
			URL   string `xml:"url"`
			Link  string `xml:"link"`
		} `xml:"image"`
		Item Items `xml:"item"`
	} `xml:"channel"`
}
