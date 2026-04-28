package backend

import mrn "github.com/YoshihideShirai/marionette"

type (
	App          = mrn.App
	Context      = mrn.Context
	Handler      = mrn.Handler
	FlashLevel   = mrn.FlashLevel
	FlashMessage = mrn.FlashMessage
)

const (
	FlashSuccess = mrn.FlashSuccess
	FlashError   = mrn.FlashError
	FlashInfo    = mrn.FlashInfo
	FlashWarn    = mrn.FlashWarn
)

var New = mrn.New
