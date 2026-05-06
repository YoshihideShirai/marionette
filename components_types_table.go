package marionette

import mf "github.com/YoshihideShirai/marionette/frontend"

// このファイルはTableコンポーネントのProps/DTO型と補助関数を定義する。
// テーブル表示に関する型・変換ロジックをここに集約する。

type TableColumn = mf.TableColumn
type TableComponentRow = mf.TableComponentRow
type TableProps = mf.TableProps

func TableRowValues(values ...any) TableComponentRow { return mf.TableRowValues(values...) }
