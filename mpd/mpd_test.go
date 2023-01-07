package mpd_test

import (
	"bytes"
	"io/fs"
	"os"
	"testing"

	"github.com/Eyevinn/dash-mpd/xml"
	"github.com/google/go-cmp/cmp"

	m "github.com/Eyevinn/dash-mpd/mpd"

	"aqwari.net/xml/xmltree"
	"github.com/stretchr/testify/require"
)

func TestDecodeEncodeMPDs(t *testing.T) {
	testDirs := []string{"testdata/go-dash-fixtures", "testdata/schema-mpds"}
	for _, testDir := range testDirs {
		fsys := os.DirFS(testDir)
		mpdFiles, err := fs.Glob(fsys, "*.mpd")
		require.NoError(t, err)
		for _, fName := range mpdFiles {
			td, err := fs.ReadFile(fsys, fName)
			require.NoError(t, err)
			mpd := m.MPD{}
			err = xml.Unmarshal(td, &mpd)
			if fName == "invalid.mpd" {
				require.Error(t, err, "")
				continue
			}
			require.NoError(t, err, fName)
			out, err := xml.MarshalIndent(mpd, "", "  ")
			require.NoError(t, err)
			inTree, err := xmltree.Parse(td)
			require.NoError(t, err)
			outTree, err := xmltree.Parse(out)
			require.NoError(t, err)
			if !xmltree.Equal(inTree, outTree) {
				inBuf := bytes.Buffer{}
				err = xmltree.Encode(&inBuf, inTree)
				require.NoError(t, err)
				outBuf := bytes.Buffer{}
				err = xmltree.Encode(&outBuf, outTree)
				require.NoError(t, err)
				d := cmp.Diff(inBuf.String(), outBuf.String())
				// Note. There is no canonicalization and there may
				// be comments in the input, so the diff is not minimal.
				t.Errorf("non-minimal diff for mpd %s:\n%s\n", fName, d[:400])
				ofh, err := os.Create(fName)
				require.NoError(t, err)
				ofh.Write(out)
				ofh.Close()
			}
		}
	}
}
