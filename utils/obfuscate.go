package utils

import "regexp"

var urlRegexp = regexp.MustCompile(`[a-z]+://([^\s]*)`)

// ObfuscateURLs obfuscates URLs in given string.
func ObfuscateURLs(src string) string {
	var offset int
	for _, sub := range urlRegexp.FindAllStringSubmatchIndex(src, -1) {
		if len(sub) < 4 {
			continue
		}
		src = src[:sub[2]-offset] + "***" + src[sub[3]-offset:]
		offset += sub[3] - sub[2] - 3
	}
	return src
}
