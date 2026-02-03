package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"passvault-fyne/internal/database"
	"passvault-fyne/pkg/utils"
)

type FieldEditor struct {
	Container *fyne.Container
	Fields    []*FieldRow
}

type FieldRow struct {
	ID         string
	KeyEntry   *widget.Entry
	ValueEntry *widget.Entry
	Sensitive  *widget.Check
	RemoveBtn  *widget.Button
	Row        *fyne.Container
}

func NewFieldEditor() *FieldEditor {
	return &FieldEditor{
		Container: container.NewVBox(),
		Fields:    []*FieldRow{},
	}
}

func (fe *FieldEditor) AddField(field database.Field) {
	if field.ID == "" {
		field.ID = utils.NewUUID()
	}
	keyEntry := widget.NewEntry()
	keyEntry.SetPlaceHolder("Field Name")
	keyEntry.Text = field.Key
	keyEntry.Wrapping = fyne.TextWrapOff

	valueEntry := widget.NewEntry()
	valueEntry.SetPlaceHolder("Value")
	valueEntry.Text = string(field.Value)
	valueEntry.Wrapping = fyne.TextWrapOff

	sensitiveCheck := widget.NewCheck("Sensitive", nil)
	sensitiveCheck.Checked = field.IsSensitive

	row := &FieldRow{
		ID:       field.ID,
		KeyEntry: keyEntry,

		ValueEntry: valueEntry,
		Sensitive:  sensitiveCheck,
	}

	removeBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		fe.RemoveField(row)
	})
	row.RemoveBtn = removeBtn

	keyWrap := container.NewMax(container.NewPadded(keyEntry))
	valueWrap := container.NewMax(container.NewPadded(valueEntry))
	rowContent := container.NewGridWithColumns(4,
		container.NewMax(keyWrap),
		container.NewMax(valueWrap),
		container.NewHBox(sensitiveCheck, layout.NewSpacer()),
		layout.NewSpacer(),
	)
	finalRow := container.NewBorder(nil, nil, nil, removeBtn, rowContent)
	row.Row = finalRow

	fe.Fields = append(fe.Fields, row)
	fe.Container.Add(finalRow)
	fe.Container.Refresh()
}

func (fe *FieldEditor) RemoveField(row *FieldRow) {
	for i, f := range fe.Fields {
		if f == row {
			fe.Fields = append(fe.Fields[:i], fe.Fields[i+1:]...)
			break
		}
	}
	fe.Container.Remove(row.Row)
	fe.Container.Refresh()
}

func (fe *FieldEditor) GetFields() []database.Field {
	var fields []database.Field
	for _, row := range fe.Fields {
		id := row.ID
		if id == "" {
			id = utils.NewUUID()
		}
		fields = append(fields, database.Field{
			ID:          id,
			Key:         row.KeyEntry.Text,
			Value:       []byte(row.ValueEntry.Text),
			IsSensitive: row.Sensitive.Checked,
		})
	}
	return fields
}
