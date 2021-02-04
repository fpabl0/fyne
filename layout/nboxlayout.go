package layout

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// CrossAlignment defines cross axis alignment type.
type CrossAlignment int

// Cross Axis alignment options.
const (
	CrossAlignmentStart CrossAlignment = iota
	CrossAlignmentEnd
	CrossAlignmentCenter
	CrossAlignmentBaseline
	CrossAlignmentStretch
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*nboxLayout)(nil)

type nboxLayout struct {
	expanded       bool
	horizontal     bool
	crossAlignment CrossAlignment
}

// ===============================================================
// Constructors
// ===============================================================

// NewNHBoxLayout returns a horizontal box layout for stacking a number of child
// canvas objects or widgets left to right.
func NewNHBoxLayout() fyne.Layout {
	return &nboxLayout{false, true, CrossAlignmentStretch}
}

// NewHBoxAlignedLayout ...
func NewHBoxAlignedLayout(crossAlignment CrossAlignment) fyne.Layout {
	return &nboxLayout{false, true, crossAlignment}
}

// NewNVBoxLayout returns a vertical box layout for stacking a number of child
// canvas objects or widgets top to bottom.
func NewNVBoxLayout() fyne.Layout {
	return &nboxLayout{false, false, CrossAlignmentStretch}
}

// NewVBoxAlignedLayout ...
func NewVBoxAlignedLayout(crossAlignment CrossAlignment) fyne.Layout {
	return &nboxLayout{false, false, crossAlignment}
}

// NewHBoxExpandedLayout ...
func NewHBoxExpandedLayout() fyne.Layout {
	return &nboxLayout{true, true, CrossAlignmentStretch}
}

// NewHBoxExpandedAlignedLayout ...
func NewHBoxExpandedAlignedLayout(crossAlignment CrossAlignment) fyne.Layout {
	return &nboxLayout{true, true, crossAlignment}
}

// NewVBoxExpandedLayout ...
func NewVBoxExpandedLayout() fyne.Layout {
	return &nboxLayout{true, false, CrossAlignmentStretch}
}

// NewVBoxExpandedAlignedLayout ...
func NewVBoxExpandedAlignedLayout(crossAlignment CrossAlignment) fyne.Layout {
	return &nboxLayout{true, false, crossAlignment}
}

// ===============================================================
// Implementation
// ===============================================================

func (g *nboxLayout) isSpacer(obj fyne.CanvasObject) bool {
	// invisible spacers don't impact layout
	if !obj.Visible() {
		return false
	}

	if g.horizontal {
		return isHorizontalSpacer(obj)
	}
	return isVerticalSpacer(obj)
}

func (g *nboxLayout) getMainSize(obj fyne.CanvasObject, min bool) float32 {
	minSize := obj.MinSize()
	curSize := obj.Size()
	if g.horizontal {
		if min {
			return minSize.Width
		}
		return fyne.Max(curSize.Width, minSize.Width)
	}
	if min {
		return minSize.Height
	}
	return fyne.Max(curSize.Height, minSize.Height)
}

func (g *nboxLayout) getCrossSize(obj fyne.CanvasObject, min bool) float32 {
	minSize := obj.MinSize()
	curSize := obj.Size()
	if g.horizontal {
		if min {
			return minSize.Height
		}
		return fyne.Max(curSize.Height, minSize.Height)
	}
	if min {
		return minSize.Width
	}
	return fyne.Max(curSize.Width, minSize.Width)
}

func (g *nboxLayout) childBoxLayoutContainerExpand(obj fyne.CanvasObject) (mainExpand bool, crossFull bool) {
	cont, ok := obj.(*fyne.Container)
	if !ok {
		return
	}
	bl, ok := cont.Layout.(*nboxLayout)
	if !ok {
		return
	}
	if !bl.expanded {
		return
	}
	// if a HBox has as a child a ExpandedHBox then it needs main axis expand.
	// if a VBox has as a child a ExpandedVBox then it needs main axis expand.
	if (g.horizontal && bl.horizontal) ||
		(!g.horizontal && !bl.horizontal) {
		mainExpand = true
		return
	}
	// if a HBox has as a child an ExpandedVBox then it needs cross axis expand.
	// if a VBox has as a child an ExpandedHBox then it needs cross axis expand
	if (g.horizontal && !bl.horizontal) ||
		(!g.horizontal && bl.horizontal) {
		crossFull = true
		return
	}
	return
}

// Layout is called to pack all child objects into a specified size.
// For a VBoxLayout this will pack objects into a single column where each item
// is full width but the height is the minimum required.
// Any spacers added will pad the view, sharing the space if there are two or more.
func (g *nboxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	maxMainSize := size.Width
	maxCrossSize := size.Height
	if !g.horizontal {
		maxMainSize = size.Height
		maxCrossSize = size.Width
	}
	spacers := 0
	flexCount := 0
	fullCrossSizeChildren := make([]fyne.CanvasObject, 0)
	minSpacePerFlex := float32(0)
	crossSize := float32(0)
	allocatedSize := float32(0) // Sum of the sizes of the non-flexible children.

	// -- Resizing non-flexible children
	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		if g.isSpacer(child) {
			spacers++
			continue
		}
		// -- Get child minimum size
		childMainSize := g.getMainSize(child, true)
		childCrossSize := g.getCrossSize(child, true)
		// if cross alignment is stretch, then full cross size
		if g.crossAlignment == CrossAlignmentStretch {
			childCrossSize = maxCrossSize
		}
		// calculate the max cross size among all the children
		crossSize = fyne.Max(crossSize, childCrossSize)

		// -- Check if we have flexible objects
		mainExpand, crossFull := g.childBoxLayoutContainerExpand(child)
		if mainExpand || g.expanded {
			flexCount++
			minSpacePerFlex = fyne.Max(minSpacePerFlex, childMainSize)
			continue
		}

		// accumulate the main axis size
		allocatedSize += childMainSize
		// if cross full size is required, then save it for resizing later,
		// when we already know the cross size.
		if crossFull {
			fullCrossSizeChildren = append(fullCrossSizeChildren, child)
			continue
		}
		// Resize non-flexible children
		if g.horizontal {
			child.Resize(fyne.NewSize(childMainSize, childCrossSize))
		} else {
			child.Resize(fyne.NewSize(childCrossSize, childMainSize))
		}
	}

	// -- Resizing flexible children
	freeSpace := maxMainSize - allocatedSize - (theme.Padding() * float32(len(objects)-spacers-1))
	spacePerFlex := float32(0)
	if g.horizontal && !g.expanded {
		fmt.Printf("free space: %.0f, minSpacerPerFlex: %.0f, flexCount: %d, spacers: %d\n", freeSpace, minSpacePerFlex, flexCount, spacers)
	}
	// TODO should we remove spacers when min size?? If we remove spacers, then minSpacePerFlex variable
	// should be removed too (as it only has sense for this)
	// if freeSpace <= (minSpacePerFlex * float32(flexCount)) {
	// 	spacers = 0
	// }
	if (flexCount + spacers) > 0 {
		spacePerFlex = freeSpace / float32(flexCount+spacers)
	}
	if g.horizontal && !g.expanded {
		fmt.Printf("spacePerFlex: %.0f\n\n", spacePerFlex)
	}
	if flexCount > 0 {
		for _, child := range objects {
			if !child.Visible() {
				continue
			}
			if g.isSpacer(child) {
				continue
			}
			// Check if we need to do some flexible children resizing.
			mainExpand, _ := g.childBoxLayoutContainerExpand(child)
			if !mainExpand && !g.expanded {
				continue
			}
			childMainSize := spacePerFlex
			childCrossSize := g.getCrossSize(child, true)
			if g.crossAlignment == CrossAlignmentStretch {
				childCrossSize = maxCrossSize
			}
			if g.horizontal {
				child.Resize(fyne.NewSize(childMainSize, childCrossSize))
			} else {
				child.Resize(fyne.NewSize(childCrossSize, childMainSize))
			}
			allocatedSize += g.getMainSize(child, false)
		}
	}

	// -- Resizing full cross size children
	for _, child := range fullCrossSizeChildren {
		if g.horizontal {
			child.Resize(fyne.NewSize(g.getMainSize(child, true), crossSize))
		} else {
			child.Resize(fyne.NewSize(crossSize, g.getMainSize(child, true)))
		}
	}

	// -- Position children
	mainPos := float32(0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		if spacers > 0 && g.isSpacer(child) {
			mainPos += spacePerFlex
			continue
		}

		crossPos := float32(0)
		switch g.crossAlignment {
		case CrossAlignmentStart:
			crossPos = 0
		case CrossAlignmentEnd:
			crossPos = crossSize - g.getCrossSize(child, false)
		case CrossAlignmentCenter:
			crossPos = (crossSize - g.getCrossSize(child, false)) / 2
		case CrossAlignmentBaseline:
			crossPos = 0
			if g.horizontal {
				// TODO
			}
		}

		if g.horizontal {
			child.Move(fyne.NewPos(mainPos, crossPos))
		} else {
			child.Move(fyne.NewPos(crossPos, mainPos))
		}

		mainPos += theme.Padding() + g.getMainSize(child, false)
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a BoxLayout this is the width of the widest item and the height is
// the sum of of all children combined with padding between each.
func (g *nboxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	mainSize, crossSize := float32(0), float32(0)
	flexCount := 0
	spacers := 0
	minSpacePerFlex := float32(0)
	addPadding := false
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if g.isSpacer(child) {
			spacers++
			continue
		}

		mainExpand, _ := g.childBoxLayoutContainerExpand(child)
		if mainExpand || g.expanded {
			flexCount++
			minSpacePerFlex = fyne.Max(minSpacePerFlex, g.getMainSize(child, true))
		} else {
			mainSize += g.getMainSize(child, true)
		}
		crossSize = fyne.Max(crossSize, g.getCrossSize(child, true))
		if addPadding {
			mainSize += theme.Padding()
		}

		addPadding = true
	}

	mainSize += minSpacePerFlex * float32(flexCount+spacers)

	if g.horizontal {
		return fyne.NewSize(mainSize, crossSize)
	}
	return fyne.NewSize(crossSize, mainSize)
}
