package glfw

import (
	"fyne.io/fyne/v2"
)

func buildMenuOverlay(menus *fyne.MainMenu, c fyne.Canvas) fyne.CanvasObject {
	if len(menus.Items) == 0 {
		fyne.LogError("Main menu must have at least one child menu", nil)
		return nil
	}

	menus = addMissingQuit(menus)
	return NewMenuBar(menus, c)
}

func addMissingQuit(menus *fyne.MainMenu) *fyne.MainMenu {
	var lastItem *fyne.MenuItem
	if len(menus.Items[0].Items) > 0 {
		lastItem = menus.Items[0].Items[len(menus.Items[0].Items)-1]
		if lastItem.Label == "Quit" {
			lastItem.IsQuit = true
		}
	}
	if lastItem == nil || !lastItem.IsQuit { // make sure the first menu always has a quit option
		quitItem := fyne.NewMenuItem("Quit", nil)
		menus.Items[0].Items = append(menus.Items[0].Items, fyne.NewMenuItemSeparator(), quitItem)
	}
	for _, item := range menus.Items[0].Items {
		if item.IsQuit && item.Action == nil {
			item.Action = func() {
				fyne.CurrentApp().Quit()
			}
		}
	}
	return menus
}
