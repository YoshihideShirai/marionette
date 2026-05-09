package frontend

import shared "github.com/YoshihideShirai/marionette/frontend/shared"

// このファイルは基本的なコンポーネントProps/DTO型を定義する。
// 新しい汎用UIコンポーネントの型はfrontend/sharedを正として公開する。

// ComponentProps defines shared style knobs for template components.
type ComponentProps = shared.ComponentProps

type LinkProps = shared.LinkProps

type IconButtonProps = shared.IconButtonProps

type LoginButtonProps = shared.LoginButtonProps

type SelectOption = shared.SelectOption

type ModalProps = shared.ModalProps

type FormFieldProps = shared.FormFieldProps

type FormProps = shared.FormProps

type ActionFormProps = shared.ActionFormProps

type StreamTriggerProps = shared.StreamTriggerProps

type InputOptions = shared.InputOptions

type TextareaOptions = shared.TextareaOptions

type EmptyStateProps = shared.EmptyStateProps

type AlertProps = shared.AlertProps

type ToastProps = shared.ToastProps

type SkeletonProps = shared.SkeletonProps

type ProgressProps = shared.ProgressProps

type VariantToken = shared.VariantToken

const (
	VariantDefault   = shared.VariantDefault
	VariantPrimary   = shared.VariantPrimary
	VariantSecondary = shared.VariantSecondary
	VariantAccent    = shared.VariantAccent
	VariantNeutral   = shared.VariantNeutral
	VariantInfo      = shared.VariantInfo
	VariantSuccess   = shared.VariantSuccess
	VariantWarning   = shared.VariantWarning
	VariantError     = shared.VariantError
	VariantGhost     = shared.VariantGhost
	VariantOutline   = shared.VariantOutline
	VariantDash      = shared.VariantDash
	VariantSoft      = shared.VariantSoft
	VariantLink      = shared.VariantLink
)

type SizeToken = shared.SizeToken

const (
	SizeXS = shared.SizeXS
	SizeSM = shared.SizeSM
	SizeMD = shared.SizeMD
	SizeLG = shared.SizeLG
	SizeXL = shared.SizeXL
)

type ButtonVariantProps = shared.ButtonVariantProps

type InputVariantProps = shared.InputVariantProps

type SelectVariantProps = shared.InputVariantProps
type BadgeVariantProps = shared.InputVariantProps
type AlertVariantProps = shared.InputVariantProps
type TextareaVariantProps = shared.InputVariantProps
type CheckboxVariantProps = shared.InputVariantProps
type RadioVariantProps = shared.InputVariantProps
type SwitchVariantProps = shared.InputVariantProps
type ProgressVariantProps = shared.InputVariantProps
type TabsVariantProps = shared.InputVariantProps
type CardVariantProps = shared.InputVariantProps
type ModalVariantProps = shared.InputVariantProps
type DrawerVariantProps = shared.InputVariantProps
type PaginationVariantProps = shared.InputVariantProps
