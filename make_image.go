package main

import (
	"github.com/LiaungYip/glyphs"
	"github.com/chai2010/webp"
	"image"
	"io/ioutil"
	"log"
)

func tempImageDir() string {
	dir, err := ioutil.TempDir("", "ingress_glyph_bot")
	check(err)
	log.Printf("Creating temporary image directory %s", dir)
	return dir
}

func makeImage(glyphNames []string) image.Image {
	s := glyphs.DefaultSettings(200)
	img := glyphs.DrawGlyphSequence(glyphNames, s)
	return img
}

func encodeWebp(img image.Image, dir string) string {
	f, err := ioutil.TempFile(dir, "glyph_")
	check(err)
	webp.Encode(f, img, nil)
	f.Close()
	return f.Name()
}
