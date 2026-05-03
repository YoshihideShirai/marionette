package goexamples

import (
	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func RegisterCarouselItemExample(app *mb.App) {
	app.Page("/carousel-item", func(ctx *mb.Context) mf.Node {
		return mf.Carousel(mf.CarouselItem("slide1", mf.Image(mf.ImageProps{Src: "https://placehold.co/320x120", Alt: "slide", Fit: "cover"})))
	})
}
