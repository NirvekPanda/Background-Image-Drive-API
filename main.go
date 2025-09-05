package main

import (
	"fmt"
)

// create a api for getting and sending images to google drive

// save meta data for all images and their locations/titles

// data structs needed:
// image_request
// image_response
// image_meta_data
// location data
// lat/long
// title
// location name

// services :
// GET getCurrentImage
// POST uploadImage
// GET imageCount

// subfunctions

// connect to google drive
// connect to location api/package

// extract location metadata
// get from location name
// get from coordinate data
// set location metadata
// based on coords or loc name, set alt

func main() {
	s := "gopher"
	fmt.Printf("Hello and welcome, %s!\n", s)

	for i := 1; i <= 5; i++ {
		fmt.Println("i =", 100/i)
	}
}
