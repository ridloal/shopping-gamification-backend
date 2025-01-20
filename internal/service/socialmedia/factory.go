package socialmedia

import "fmt"

func GetVerifier(platform string) (SocialMediaVerifier, error) {
	switch platform {
	case "instagram":
		return NewInstagramVerifier(), nil
	case "tiktok":
		return NewTikTokVerifier(), nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}
