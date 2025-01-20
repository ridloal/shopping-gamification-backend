package socialmedia

type SocialMediaVerifier interface {
	VerifyLike(username, postID string) (bool, error)
	VerifyComment(username, postID string) (bool, error)
	VerifyShare(username, postID string) (bool, error)
	VerifyFollow(username, targetUsername string) (bool, error)
}
