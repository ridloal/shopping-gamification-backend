package socialmedia

import (
	"net/http"
	"os"
)

type InstagramVerifier struct {
	accessToken string
	apiVersion  string
	client      *http.Client
}

func NewInstagramVerifier() *InstagramVerifier {
	return &InstagramVerifier{
		accessToken: os.Getenv("INSTAGRAM_ACCESS_TOKEN"),
		apiVersion:  "v12.0",
		client:      &http.Client{},
	}
}

func (iv *InstagramVerifier) VerifyLike(username, postID string) (bool, error) {
	// url := fmt.Sprintf("https://graph.instagram.com/%s/%s/likes?access_token=%s", iv.apiVersion, postID, iv.accessToken)

	// Implementation of Instagram API call
	return false, nil
}

func (iv *InstagramVerifier) VerifyComment(username, postID string) (bool, error) {

	// Implementation of Instagram API call
	return false, nil
}

func (iv *InstagramVerifier) VerifyShare(username, postID string) (bool, error) {

	// Implementation of Instagram API call
	return false, nil
}

func (iv *InstagramVerifier) VerifyFollow(username, targetUsername string) (bool, error) {

	// Implementation of Instagram API call
	return false, nil
}
