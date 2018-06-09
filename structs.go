package main

type catImage struct {
	URL string `xml:"data>images>image>url"`
}
