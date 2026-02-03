package ui

import (
	"bytes"
	"encoding/json"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (a *App) showSecretDetails(id string) {
	secret, err := a.DB.GetSecret(id)
	if err != nil {
		dialog.ShowError(err, a.MainWindow)
		return
	}

	if secret == nil {
		return
	}

	nameLabel := widget.NewLabelWithStyle(secret.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	fieldsContainer := container.NewVBox()
	var contentArea fyne.CanvasObject
	if secret.Format == "json" || secret.Format == "text" {
		contentText := secret.Content
		if secret.Format == "json" {
			var pretty bytes.Buffer
			if err := json.Indent(&pretty, []byte(secret.Content), "", "  "); err == nil {
				contentText = pretty.String()
			}
		}
		contentEntry := widget.NewMultiLineEntry()
		contentEntry.SetText(contentText)
		contentEntry.Wrapping = fyne.TextWrapWord
		contentEntry.Scroll = container.ScrollVerticalOnly
		contentEntry.SetMinRowsVisible(18)
		copyBtn := widget.NewButton("Copy", func() {
			selected := contentEntry.SelectedText()
			if selected == "" {
				selected = contentText
			}
			a.Clipboard.CopyWithAutoClear(selected, 30*time.Second)
		})
		headRow := container.NewHBox(widget.NewLabel("Content:"), layout.NewSpacer(), copyBtn)
		split := container.NewVSplit(contentEntry, layout.NewSpacer())
		split.SetOffset(0.8)
		prefKey := "content_split_offset"
		prefs := a.FyneApp.Preferences()
		stored := prefs.FloatWithFallback(prefKey, 0.8)
		if stored < 0.4 {
			stored = 0.4
		} else if stored > 0.95 {
			stored = 0.95
		}
		split.SetOffset(stored)
		if a.contentSplitStop != nil {
			close(a.contentSplitStop)
		}
		a.contentSplitStop = make(chan struct{})
		go func(stop <-chan struct{}) {
			ticker := time.NewTicker(250 * time.Millisecond)
			defer ticker.Stop()
			last := split.Offset
			for {
				select {
				case <-ticker.C:
					if split.Offset != last {
						last = split.Offset
						prefs.SetFloat(prefKey, last)
					}
				case <-stop:
					return
				}
			}
		}(a.contentSplitStop)
		contentArea = container.NewBorder(headRow, nil, nil, nil, split)
	} else {
		for _, field := range secret.Fields {
			valStr := string(field.Value)
			if field.IsSensitive {
				valStr = "********"
			}

			valLabel := widget.NewLabel(valStr)
			keyLabel := widget.NewLabelWithStyle(field.Key+":", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})

			copyBtn := widget.NewButton("Copy", func() {
				a.Clipboard.CopyWithAutoClear(string(field.Value), 30*time.Second)
			})

			row := container.NewHBox(keyLabel, valLabel, layout.NewSpacer(), copyBtn)
			fieldsContainer.Add(row)
		}
		contentArea = container.NewVScroll(fieldsContainer)
	}

	editBtn := widget.NewButton("Edit", func() {
		a.ShowEditSecretDialog(secret)
	})

	deleteBtn := widget.NewButton("Delete", func() {
		dialog.ShowConfirm("Delete Secret", "Are you sure?", func(ok bool) {
			if ok {
				if err := a.DB.DeleteSecret(secret.ID); err != nil {
					dialog.ShowError(err, a.MainWindow)
				} else {
					a.refreshSecretList()
					empty := widget.NewLabel("Select a secret")

					split := a.MainWindow.Content().(*container.Split)
					split.Trailing = empty
					split.Refresh()
				}
			}
		}, a.MainWindow)
	})

	topBar := container.NewBorder(nil, nil, nameLabel, container.NewHBox(editBtn, deleteBtn), nil)

	if contentArea == nil {
		contentArea = container.NewVScroll(fieldsContainer)
	}
	contentWrap := container.NewBorder(widget.NewSeparator(), nil, nil, nil, contentArea)
	detailView := container.NewBorder(topBar, nil, nil, nil, contentWrap)

	split := a.MainWindow.Content().(*container.Split)
	split.Trailing = container.NewPadded(detailView)
	split.Refresh()
}
