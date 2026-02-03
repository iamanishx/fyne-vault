package ui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"passvault-fyne/internal/database"
	"passvault-fyne/pkg/utils"
)

func (a *App) ShowAddSecretDialog() {
	a.showSecretDialog(nil)
}

func (a *App) ShowEditSecretDialog(secret *database.SecretEntry) {
	a.showSecretDialog(secret)
}

func (a *App) showSecretDialog(existing *database.SecretEntry) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Name")

	fieldEditor := NewFieldEditor()

	if existing != nil {
		nameEntry.Text = existing.Name
		for _, f := range existing.Fields {
			fieldEditor.AddField(f)
		}
	} else {
		fieldEditor.AddField(database.Field{Key: "Username", Value: []byte(""), IsSensitive: false})
		fieldEditor.AddField(database.Field{Key: "Password", Value: []byte(""), IsSensitive: true})
	}

	addFieldBtn := widget.NewButton("Add Field", func() {
		fieldEditor.AddField(database.Field{ID: utils.NewUUID()})
	})
	addFieldBtn.Importance = widget.LowImportance

	nameRow := container.NewBorder(nil, nil, widget.NewLabel("Name"), nil, nameEntry)
	fieldsHeader := container.NewBorder(nil, nil, widget.NewLabel("Fields"), nil, addFieldBtn)
	fieldsBox := container.NewVBox(fieldEditor.Container)
	fieldsCard := canvas.NewRectangle(color.NRGBA{R: 22, G: 22, B: 26, A: 255})
	fieldsCard.CornerRadius = 10
	fieldsCard.StrokeColor = color.NRGBA{R: 45, G: 45, B: 52, A: 200}
	fieldsCard.StrokeWidth = 1
	fieldsPanel := container.NewMax(fieldsCard, container.NewPadded(fieldsBox))

	content := container.NewVBox(
		widget.NewLabel("Secret Details"),
		nameRow,
		widget.NewSeparator(),
		fieldsHeader,
		fieldsPanel,
	)

	card := canvas.NewRectangle(color.NRGBA{R: 18, G: 18, B: 22, A: 255})
	card.CornerRadius = 12
	card.StrokeColor = color.NRGBA{R: 40, G: 40, B: 46, A: 200}
	card.StrokeWidth = 1
	content = container.NewPadded(content)
	content = container.NewMax(card, content)

	title := "Add Secret"
	if existing != nil {
		title = "Edit Secret"
	}

	d := dialog.NewCustomConfirm(title, "Save", "Cancel", content, func(ok bool) {
		if !ok {
			return
		}

		id := utils.NewUUID()
		createdAt := time.Now()
		if existing != nil {
			id = existing.ID
			createdAt = existing.CreatedAt
		}

		secret := &database.SecretEntry{
			ID:        id,
			Name:      nameEntry.Text,
			CreatedAt: createdAt,
			UpdatedAt: time.Now(),
			Fields:    fieldEditor.GetFields(),
		}

		if err := a.DB.SaveSecret(secret); err != nil {
			dialog.ShowError(err, a.MainWindow)
		} else {
			a.refreshSecretList()
			if existing != nil {
				a.showSecretDetails(secret.ID)
			}
		}
	}, a.MainWindow)

	d.Resize(fyne.NewSize(520, 620))
	d.Show()
}
