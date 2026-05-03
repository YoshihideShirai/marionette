package marionette

import (
	rdf "github.com/rocketlaunchr/dataframe-go"
)

func InputWithOptions(name, value string, options InputOptions) Node {
	return inputWithOptionsComponent(name, value, options)
}
func DataFrameComponent(df *rdf.DataFrame, props TableProps) Node { return DataFrame(df, props) }
