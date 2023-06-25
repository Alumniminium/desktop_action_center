package View

import (
	"fmt"

	"github.com/actionCenter/Model"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type NotificationWidget struct {
	container *gtk.Box
	id        int
}
type NotificationList struct {
	container     *gtk.ScrolledWindow
	listBox       *gtk.ListBox
	notifications []NotificationWidget
}

func (app *ActionCenterUI) createNotificationComponent() (*gtk.Box, error) {
	scrollBox, _ := gtk.ScrolledWindowNew(nil, nil)
	container, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)

	scrollBox.SetVExpand(true)
	scrollBox.SetHExpand(false)

	label, err := gtk.LabelNew("Notifications")
	if err != nil {
		return nil, err
	}
	container.Add(label)

	clearBtn, err := gtk.ButtonNewWithLabel("Clear")
	if err != nil {
		return nil, err
	}
	clearBtn.Connect("clicked", func() {
		app.clearNotification()
	})
	container.Add(clearBtn)

	listBox, _ := gtk.ListBoxNew()
	style, err := listBox.GetStyleContext()
	if err != nil {
		return nil, err
	}
	style.AddClass("notification-container")
	style.AddProvider(app.componentStyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	listBox.SetSelectionMode(gtk.SELECTION_NONE)

	nlist := NotificationList{
		container: scrollBox,
		listBox:   listBox,
	}
	app.notifications = nlist
	scrollBox.Add(listBox)
	container.Add(scrollBox)
	return container, nil
}

func (app *ActionCenterUI) AddNotification(n Model.Notification) error {
	widget := NotificationWidget{}
	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 20)
	elementWidth := app.notifications.listBox.GetAllocatedWidth() - ICON_SIZE - HORIZONTAL_SPACING

	var icon *gtk.Image = nil

	if _, ok := n.Hints["image-data"]; ok {
		width := n.Hints["image-data"].Value().([]interface{})[0].(int32)
		height := n.Hints["image-data"].Value().([]interface{})[1].(int32)
		rowStride := n.Hints["image-data"].Value().([]interface{})[2].(int32)
		hasAlpha := n.Hints["image-data"].Value().([]interface{})[3].(bool)
		bitsPerSample := n.Hints["image-data"].Value().([]interface{})[4].(int32)

		img := n.Hints["image-data"].Value().([]interface{})[6].([]byte)
		pixbuf, err := gdk.PixbufNewFromData(img, gdk.COLORSPACE_RGB, hasAlpha, int(bitsPerSample), int(width), int(height), int(rowStride))
		icon, err = gtk.ImageNewFromPixbuf(pixbuf)
		if err != nil {
			fmt.Println(err)
		}
	} else if customImagePath, ok := n.Hints["image-path"].Value().(string); ok {
		icon, err = gtk.ImageNewFromFile(customImagePath)
	} else if n.AppIcon != "" {
		icon, err = gtk.ImageNewFromIconName(n.AppIcon, gtk.ICON_SIZE_LARGE_TOOLBAR)
	} else {
		icon, err = gtk.ImageNewFromIconName("gtk-dialog-info", gtk.ICON_SIZE_LARGE_TOOLBAR)
	}

	if icon != nil {
		resize(icon)
		hbox.PackStart(icon, false, false, 0)
	}

	summaryLabel, err := gtk.LabelNew(n.Summary)
	stylectx, err := summaryLabel.GetStyleContext()
	stylectx.AddClass("notification-summary")
	stylectx.AddProvider(app.componentStyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	summaryLabel.SetHAlign(gtk.ALIGN_START)
	summaryLabel.SetLineWrap(true)
	summaryLabel.SetMaxWidthChars(1)
	summaryLabel.SetSizeRequest(elementWidth, -1)
	summaryLabel.SetXAlign(0)

	bodyLabel, err := gtk.LabelNew(n.Body)
	stylectx, err = bodyLabel.GetStyleContext()
	stylectx.AddClass("notification-body")
	stylectx.AddProvider(app.componentStyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	bodyLabel.SetLineWrap(true)
	bodyLabel.SetMaxWidthChars(1)
	bodyLabel.SetSizeRequest(elementWidth, -1)
	bodyLabel.SetHAlign(gtk.ALIGN_START)
	bodyLabel.SetXAlign(0)

	vbox.PackStart(summaryLabel, false, false, 0)
	vbox.PackStart(bodyLabel, false, false, 0)
	hbox.PackStart(vbox, true, true, 0)

	widget.container = hbox

	row, err := gtk.ListBoxRowNew()
	row.Add(widget.container)

	app.notifications.listBox.Insert(row, 0)
	return err
}

func (app *ActionCenterUI) clearNotification() {
	for app.notifications.listBox.GetChildren().Length() > 0 {
		app.notifications.listBox.Remove(app.notifications.listBox.GetRowAtIndex(0))
	}
}

func resize(icon *gtk.Image) {
	pixbuf := icon.GetPixbuf()
	if pixbuf == nil {
		theme, _ := gtk.IconThemeGetDefault()
		iconName, _ := icon.GetIconName()
		pixbuf, _ = theme.LoadIconForScale(iconName, ICON_SIZE, 1, gtk.ICON_LOOKUP_FORCE_SIZE)
	}
	scaledPixbuf, _ := pixbuf.ScaleSimple(ICON_SIZE, ICON_SIZE, gdk.INTERP_BILINEAR)
	icon.SetFromPixbuf(scaledPixbuf)
}
