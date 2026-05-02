package main

import (
	"os"

	mb "github.com/YoshihideShirai/marionette/backend"
	mf "github.com/YoshihideShirai/marionette/frontend"
)

func main() {
	_ = os.ErrNotExist
	_ = mb.New
	_ = mf.DataFrameFromCSV
	_ = mf.DataFrameFromTSV
}
