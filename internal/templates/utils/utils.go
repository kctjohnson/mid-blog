package utils

import "github.com/go-loremipsum/loremipsum"

var WordGenerator = newPlaceholder()

type placeholder struct {
	LoremIpsum *loremipsum.LoremIpsum
}

func newPlaceholder() *placeholder {
	return &placeholder{
		LoremIpsum: loremipsum.New(),
	}
}

func (p *placeholder) Words(count int) string {
	return p.LoremIpsum.Words(count)
}

func (p *placeholder) Sentences(count int) string {
	return p.LoremIpsum.Sentences(count)
}

func (p *placeholder) Paragraphs(count int) string {
	return p.LoremIpsum.Paragraphs(count)
}
