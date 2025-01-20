package socialmedia

import (
	"net/http"
	"os"
)

type TikTokVerifier struct {
	accessToken string
	client      *http.Client
}

func NewTikTokVerifier() *TikTokVerifier {
	return &TikTokVerifier{
		accessToken: os.Getenv("TIKTOK_ACCESS_TOKEN"),
		client:      &http.Client{},
	}
}

func (tv *TikTokVerifier) VerifyLike(username, postID string) (bool, error) {
	// url := fmt.Sprintf("https://open-api.tiktok.com/video/likes?access_token=%s&post_id=%s", tv.accessToken, postID)

	// Implementation of TikTok API call
	return false, nil
}

func (tv *TikTokVerifier) VerifyComment(username, postID string) (bool, error) {

	// Implementation of TikTok API call
	return false, nil
}

func (tv *TikTokVerifier) VerifyShare(username, postID string) (bool, error) {

	// Implementation of TikTok API call
	return false, nil
}

func (tv *TikTokVerifier) VerifyFollow(username, targetUsername string) (bool, error) {

	// Implementation of TikTok API call
	return false, nil
}
