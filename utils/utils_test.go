package utils

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"testing"

	"github.com/twsiyuan/evernote-sdk-golang/notestore"
	"github.com/twsiyuan/evernote-sdk-golang/types"
)

const (
	EvernoteEnvironment EnvironmentType = SANDBOX

	// Get API key, see https://dev.evernote.com/doc/
	EvernoteKey    string = ""
	EvernoteSecret string = ""

	// For test, you can get develper token, see https://dev.evernote.com/doc/articles/dev_tokens.php
	EvernoteAuthorToken string = ""
)

func TestMain(t *testing.T) {
	us, err := GetUserStore(EvernoteEnvironment)
	if err != nil {
		t.Fatal(err)
	}
	ns, err := GetNoteStore(us, EvernoteAuthorToken)
	if err != nil {
		t.Fatal(err)
	}
	{
		note, err := ns.GetDefaultNotebook(EvernoteAuthorToken)
		if err != nil {
			t.Fatal(err)
		}
		if note == nil {
			t.Fatal("Invalid Note")
		}
		t.Logf("Default notebook: %s", note.GetName())
	}

	var testNotebook *types.Notebook = nil
	{
		notebooks, err := ns.ListNotebooks(EvernoteAuthorToken)
		if err != nil {
			t.Fatal(err)
		}
		if notebooks == nil {
			t.Fatal("Invalid Notebooks")
		}
		for idx, notebook := range notebooks {
			t.Logf("Notebook[%d]: %s", idx, notebook.GetName())
			testNotebook = notebook
		}
	}
	var testNote *notestore.NoteMetadata = nil
	if testNotebook != nil {
		maxNotes := int32(250)

		filter := notestore.NewNoteFilter()
		filter.NotebookGuid = testNotebook.GUID

		it := true
		resultSpec := notestore.NewNotesMetadataResultSpec()
		resultSpec.IncludeTitle = &it

		notes, err := ns.FindNotesMetadata(EvernoteAuthorToken, filter, 0, maxNotes, resultSpec)
		if err != nil {
			t.Fatal(err)
		}

		for idx, note := range notes.GetNotes() {
			t.Logf("Note[%d]: %s, guid: %s", idx, note.GetTitle(), note.GetGUID())
			testNote = note
		}
	}
	if testNote != nil {
		note, err := ns.GetNote(EvernoteAuthorToken, testNote.GetGUID(), true, true, true, true)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Note Content: %s", note.GetContent())

		for idx, res := range note.GetResources() {
			t.Logf("Note Resources[%d]: %s, %d byte(s), %s", idx, res.GetMime(), res.GetData().GetSize(), res.GetGUID())

			if res.GetMime() == "image/png" {
				img, err := png.Decode(bytes.NewBuffer(res.GetData().Body))
				if err != nil {
					t.Error(err)
				} else {
					{
						f, err := os.Create(fmt.Sprintf("Z:/%s.%s", res.GetGUID(), "jpg"))
						if err != nil {
							t.Error(err)
						}

						jpeg.Encode(f, img, &jpeg.Options{95})

					}
					{
						f, err := os.Create(fmt.Sprintf("Z:/%s.%s", res.GetGUID(), "png"))
						if err != nil {
							t.Error(err)
						}

						png.Encode(f, img)
					}
				}
			}
		}
	}
}
