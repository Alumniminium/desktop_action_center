package Radio

import (
	"github.com/actionCenter/Data"
	"github.com/fhs/gompd/mpd"
	"github.com/gotk3/gotk3/gtk"
)

type RadioTab struct {
	container         *gtk.Box
	listbox           *gtk.ListBox
	directoryServerIp string
	foundStations     []Station
	mpdClient         mpd.Client
}

func (radio *RadioTab) Create() (*gtk.Box, error) {
	mpdclient, _ := mpd.Dial("tcp", "localhost:6600")
	radio.mpdClient = *mpdclient
	radio.mpdClient.Clear()
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	// Music player
	playerBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	label, _ := gtk.LabelNew("loading")
	commandBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	stopButton, _ := gtk.ButtonNewWithLabel("stop")
	voteButton, _ := gtk.ButtonNewWithLabel("vote")
	commandBox.Add(stopButton)
	commandBox.Add(voteButton)
	playerBox.Add(label)
	playerBox.Add(commandBox)
	if err != nil {
		return nil, err
	}
	// Advanced search
	advancedSearchBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	inputBox, _ := gtk.EntryNew()
	listBox, _ := gtk.ListBoxNew()
	listBox.Connect("row-selected", func() {
		selected := listBox.GetSelectedRow()
		if selected != nil {
			radio.mpdClient.Clear()
			label.SetText(radio.foundStations[selected.GetIndex()].Name)
			radio.mpdClient.Add(radio.foundStations[selected.GetIndex()].URL)
			radio.mpdClient.Play(0)
		}
	})
	inputBox.Connect("activate", func() {
		text, _ := inputBox.GetText()
		inputBox.SetText("")
		stations := radio.AdvancedStationSearch(text, "", 3)
		for listBox.GetChildren().Length() > 0 {
			listBox.Remove(listBox.GetRowAtIndex(0))
		}
		for _, s := range stations {
			radio.AddStation(s)
		}
	})
	style, _ := inputBox.GetStyleContext()
	style.AddClass("ai-inputbox")
	style.AddProvider(Data.StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	advancedSearchBox.Add(listBox)
	advancedSearchBox.Add(inputBox)

	radio.listbox = listBox
	radio.container = container
	container.Add(playerBox)
	container.Add(advancedSearchBox)

	return container, nil
}
