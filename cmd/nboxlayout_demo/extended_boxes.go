package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func newVBox(objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(layout.NewNVBoxLayout(), objects...)
}

func newVBoxAligned(crossAlignment layout.CrossAlignment, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(layout.NewVBoxAlignedLayout(crossAlignment), objects...)
}

func newVBoxExpanded(objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(layout.NewVBoxExpandedLayout(), objects...)
}

func newVBoxExpandedAligned(crossAlignment layout.CrossAlignment, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(layout.NewVBoxExpandedAlignedLayout(crossAlignment), objects...)
}

func newHBox(objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(layout.NewNHBoxLayout(), objects...)
}

func newHBoxAligned(crossAlignment layout.CrossAlignment, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(layout.NewHBoxAlignedLayout(crossAlignment), objects...)
}

func newHBoxExpanded(objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(layout.NewHBoxExpandedLayout(), objects...)
}

func newHBoxExpandedAligned(crossAlignment layout.CrossAlignment, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(layout.NewHBoxExpandedAlignedLayout(crossAlignment), objects...)
}
