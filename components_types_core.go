package marionette

import mf "github.com/YoshihideShirai/marionette/frontend"

// このファイルは基本的なコンポーネントProps/DTO型を定義する。
// 新しい汎用UIコンポーネントの型はfrontendを正として公開する。

type ComponentProps = mf.ComponentProps
type LinkProps = mf.LinkProps
type IconButtonProps = mf.IconButtonProps
type LoginButtonProps = mf.LoginButtonProps
type SelectOption = mf.SelectOption
type ModalProps = mf.ModalProps
type FormFieldProps = mf.FormFieldProps
type FormProps = mf.FormProps
type ActionFormProps = mf.ActionFormProps
type InputOptions = mf.InputOptions
type TextareaOptions = mf.TextareaOptions
type EmptyStateProps = mf.EmptyStateProps
type AlertProps = mf.AlertProps
type ToastProps = mf.ToastProps
type SkeletonProps = mf.SkeletonProps
type ProgressProps = mf.ProgressProps
type VariantToken = mf.VariantToken

const (
	VariantDefault   = mf.VariantDefault
	VariantPrimary   = mf.VariantPrimary
	VariantSecondary = mf.VariantSecondary
	VariantAccent    = mf.VariantAccent
	VariantNeutral   = mf.VariantNeutral
	VariantInfo      = mf.VariantInfo
	VariantSuccess   = mf.VariantSuccess
	VariantWarning   = mf.VariantWarning
	VariantError     = mf.VariantError
	VariantGhost     = mf.VariantGhost
	VariantOutline   = mf.VariantOutline
	VariantDash      = mf.VariantDash
	VariantSoft      = mf.VariantSoft
	VariantLink      = mf.VariantLink
)

type SizeToken = mf.SizeToken

const (
	SizeXS = mf.SizeXS
	SizeSM = mf.SizeSM
	SizeMD = mf.SizeMD
	SizeLG = mf.SizeLG
	SizeXL = mf.SizeXL
)

type ButtonVariantProps = mf.ButtonVariantProps
type InputVariantProps = mf.InputVariantProps
type SelectVariantProps = mf.SelectVariantProps
type BadgeVariantProps = mf.BadgeVariantProps
type AlertVariantProps = mf.AlertVariantProps
type TextareaVariantProps = mf.TextareaVariantProps
type CheckboxVariantProps = mf.CheckboxVariantProps
type RadioVariantProps = mf.RadioVariantProps
type SwitchVariantProps = mf.SwitchVariantProps
type ProgressVariantProps = mf.ProgressVariantProps
type TabsVariantProps = mf.TabsVariantProps
type CardVariantProps = mf.CardVariantProps
type ModalVariantProps = mf.ModalVariantProps
type DrawerVariantProps = mf.DrawerVariantProps
type PaginationVariantProps = mf.PaginationVariantProps
