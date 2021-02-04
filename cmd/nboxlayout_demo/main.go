package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func formItem(label string, proposed bool) fyne.CanvasObject {
	lbl := widget.NewLabel(label)
	lbl.Alignment = fyne.TextAlignTrailing
	objs := []fyne.CanvasObject{
		/*layout.NewSpacer(),*/ lbl, newHBoxExpanded(widget.NewEntry()),
	}
	if proposed {
		return newHBox(objs...)
	}
	return container.NewHBox(objs...)
}

func buildForm(proposed bool) fyne.CanvasObject {
	objs := []fyne.CanvasObject{
		formItem("Name:", proposed),
		formItem("Last name:", proposed),
		formItem("Card:", proposed),
		layout.NewSpacer(),
	}
	if proposed {
		return newVBox(objs...)
	}
	return container.NewVBox(objs...)
}

func main() {
	a := app.New()
	w := a.NewWindow("Hello")
	w.Resize(fyne.NewSize(500, 300))

	done := make(chan int)
	go func() {
		// for {
		// 	fmt.Println(w.Content().Size())
		// 	time.Sleep(2 * time.Second)
		// 	select {
		// 	case <-done:
		// 		return
		// 	default:
		// 	}
		// }
	}()

	w.SetContent(buildForm(true))

	// w.SetContent(container.NewVBox(
	// 	proposed(),
	// 	current(true),
	// 	current(false),
	// ))

	w.ShowAndRun()
	close(done)
}

func current(border bool) fyne.CanvasObject {
	return container.NewVBox(
		cRow(border, theme.ComputerIcon()),
		cRow(border, theme.DocumentIcon()),
		cRow(border, theme.MailComposeIcon()),
	)
}

func cRow(border bool, res fyne.Resource) fyne.CanvasObject {
	img := canvas.NewImageFromResource(res)
	img.SetMinSize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	// return container.NewGridWithColumns(2, img, widget.NewEntry())
	if border {
		return container.NewBorder(nil, nil, container.NewCenter(img), nil, widget.NewEntry())
	}
	return container.NewHBox(img, widget.NewEntry())
}

func proposed() fyne.CanvasObject {
	return newVBox(
		pRow(theme.ComputerIcon()),
		pRow(theme.DocumentIcon()),
		pRow(theme.MailComposeIcon()),
	)
}

func pRow(res fyne.Resource) fyne.CanvasObject {
	img := canvas.NewImageFromResource(res)
	img.SetMinSize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	return newHBoxAligned(layout.CrossAlignmentCenter, img, newHBoxExpanded(widget.NewEntry()))
}

// func proposed() fyne.CanvasObject {
// 	return container.NewColumn(
// 		fyne.AxisAlignment{},
// 		pRow(theme.ComputerIcon()),
// 		pRow(theme.DocumentIcon()),
// 		pRow(theme.MailComposeIcon()),
// 	)
// }

// func pRow(res fyne.Resource) fyne.CanvasObject {
// 	img := canvas.NewImageFromResource(res)
// 	img.SetMinSize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
// 	return container.NewRow(
// 		fyne.AxisAlignment{},
// 		img,
// 		newHPad(),
// 		widget.NewExpanded(widget.NewEntry()),
// 	)
// }

// func newVPad() fyne.CanvasObject {
// 	r := canvas.NewRectangle(color.Transparent)
// 	r.SetMinSize(fyne.NewSize(0, theme.Padding()))
// 	return r
// }

// func newHPad() fyne.CanvasObject {
// 	r := canvas.NewRectangle(color.Transparent)
// 	r.SetMinSize(fyne.NewSize(theme.Padding(), 0))
// 	return r
// }
