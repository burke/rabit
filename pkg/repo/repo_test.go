package repo

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

const (
	blob1Path = "testfiles/blob1"
	blob2Path = "testfiles/blob2"
)

func TestGarbage(t *testing.T) {
	dir, err := ioutil.TempDir("", "rabit-test")
	if err != nil {
		t.Fatal("tempdir")
	}
	defer os.RemoveAll(dir)

	repo := New(dir)
	repo.Init()

	blob1, err := os.Open(blob1Path)
	if err != nil {
		t.Fatal("open blob1")
	}

	if err := repo.Add(blob1, "blob1"); err != nil {
		t.Fatal("repo add")
	}
	blob1.Close()

	expectChunks(t, dir, map[string]string{
		"24662838814f422b3050a99575b29a62d8af9e0f": "8b29b56689c68dd4dd8ba20170f70247",
		"270d8cd95b5f56d0153c37c17ba9bda6de181185": "e444a13162395b5b51d31d58f4c92ead",
		"35d23ee504acbb6dfd638a44167d7bb0f182d442": "531ee74a5d30a2d28bc58e1c54614138",
	})

	expectManifests(t, dir, map[string]string{
		"blob1": `24662838814f422b3050a99575b29a62d8af9e0f
270d8cd95b5f56d0153c37c17ba9bda6de181185
35d23ee504acbb6dfd638a44167d7bb0f182d442
`,
	})

	if err := repo.Rm("blob1"); err != nil {
		t.Fatal("rm fail")
	}

	expectChunks(t, dir, map[string]string{})
	expectManifests(t, dir, map[string]string{})

	blob1, err = os.Open(blob1Path)
	if err != nil {
		t.Fatal("open blob1")
	}

	if err := repo.Add(blob1, "blob1"); err != nil {
		t.Fatal("repo add")
	}
	blob1.Close()

	blob2, err := os.Open(blob2Path)
	if err != nil {
		t.Fatal("open blob2")
	}

	if err := repo.Add(blob2, "blob2"); err != nil {
		t.Fatal("repo add")
	}
	blob2.Close()

	expectChunks(t, dir, map[string]string{
		"24662838814f422b3050a99575b29a62d8af9e0f": "8b29b56689c68dd4dd8ba20170f70247",
		"270d8cd95b5f56d0153c37c17ba9bda6de181185": "e444a13162395b5b51d31d58f4c92ead",
		"073d934ddb2b9a589c91883548a20d10b2d21ec2": "028742b698eb2ace57c2c75566222996",
		"c82241b955d9d1d1acac4ca85b05e288af8f6135": "d6e1a2efa7376ada012c85379d811df7",
		"35d23ee504acbb6dfd638a44167d7bb0f182d442": "531ee74a5d30a2d28bc58e1c54614138",
	})

	expectManifests(t, dir, map[string]string{
		"blob1": `24662838814f422b3050a99575b29a62d8af9e0f
270d8cd95b5f56d0153c37c17ba9bda6de181185
35d23ee504acbb6dfd638a44167d7bb0f182d442
`,
		"blob2": `24662838814f422b3050a99575b29a62d8af9e0f
073d934ddb2b9a589c91883548a20d10b2d21ec2
c82241b955d9d1d1acac4ca85b05e288af8f6135
`,
	})

	files, err := repo.LsFiles()
	if err != nil {
		t.Error("ls failed")
	}

	sort.Strings(files)
	if !reflect.DeepEqual(files, []string{"blob1", "blob2"}) {
		t.Fatal("ls didn't look right")
	}

	if err := os.Remove(filepath.Join(dir, "manifests", "blob1")); err != nil {
		t.Fatal("os remove fail")
	}
	if err := repo.GC(false); err != nil {
		t.Fatal("GC fail")
	}

	expectChunks(t, dir, map[string]string{
		"24662838814f422b3050a99575b29a62d8af9e0f": "8b29b56689c68dd4dd8ba20170f70247",
		"073d934ddb2b9a589c91883548a20d10b2d21ec2": "028742b698eb2ace57c2c75566222996",
		"c82241b955d9d1d1acac4ca85b05e288af8f6135": "d6e1a2efa7376ada012c85379d811df7",
	})

	expectManifests(t, dir, map[string]string{
		"blob2": `24662838814f422b3050a99575b29a62d8af9e0f
073d934ddb2b9a589c91883548a20d10b2d21ec2
c82241b955d9d1d1acac4ca85b05e288af8f6135
`,
	})

	buf := bytes.NewBuffer(nil)
	repo.CatFile("blob2", buf)
	actual := buf.Bytes()
	expected, err := ioutil.ReadFile(blob2Path)

	if bytes.Compare(actual, expected) != 0 {
		t.Error("compare failed")
	}

}

func expectChunks(t *testing.T, dir string, expects map[string]string) {
	var chunks []string

	prefixes, err := ioutil.ReadDir(filepath.Join(dir, "chunks"))
	if err != nil {
		t.Fatal("readdir prefixes")
	}

	for _, pref := range prefixes {
		cinfo, err := ioutil.ReadDir(filepath.Join(dir, "chunks", pref.Name()))
		if err != nil {
			t.Fatal("readdir chunk")
		}
		for _, fi := range cinfo {
			chunks = append(chunks, fi.Name())
		}
	}

	if len(chunks) != len(expects) {
		t.Error("unexpected number of chunks")
	}

	for _, name := range chunks {
		ehash, ok := expects[name]
		if !ok {
			t.Error("unexpected chunk")
		}
		ahash, err := fileMD5(filepath.Join(dir, "chunks", name[0:2], name))
		if err != nil {
			t.Error("md5 fail", err)
		}
		if ehash != ahash {
			t.Errorf("expected %s to be %s but was %s", name, ehash, ahash)
		}
	}

}

func expectManifests(t *testing.T, dir string, expects map[string]string) {
	var manifests []string

	mfis, err := ioutil.ReadDir(filepath.Join(dir, "manifests"))
	if err != nil {
		t.Fatal("readdir manifests")
	}

	for _, mfi := range mfis {
		manifests = append(manifests, mfi.Name())
	}

	if len(manifests) != len(expects) {
		t.Error("unexpected number of manifests")
	}

	for _, name := range manifests {
		econtents, ok := expects[name]
		if !ok {
			t.Error("unexpected manifest")
		}
		acontentsBytes, err := ioutil.ReadFile(filepath.Join(dir, "manifests", name))
		if err != nil {
			t.Error("readfile fail", err)
		}
		acontents := string(acontentsBytes)
		if econtents != acontents {
			t.Errorf("expected %s to be %s but was %s", name, econtents, acontents)
		}
	}

}

func fileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
