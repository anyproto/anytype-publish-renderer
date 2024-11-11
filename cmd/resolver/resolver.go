package resolver

import "fmt"

type SimpleAssetResolver struct {
	cdnUrl string
}

func New(cdnUrl string) SimpleAssetResolver {
	return SimpleAssetResolver{
		cdnUrl: cdnUrl,
	}
}

func (r SimpleAssetResolver) ByEmojiCode(code rune) string {
	return fmt.Sprintf("%s/emojies/%x.png", r.cdnUrl, code)
}

func (r SimpleAssetResolver) ByImgPath(imagePath string) string {
	return "/static/" + imagePath
}
