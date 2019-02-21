package ghost

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Ghost struct {
}

func NewGhost() *Ghost {
	var g Ghost
	return &g
}

func (g Ghost) Convert(pdfSrc, dest string, ppi int) (string, error) {
	output := filepath.Join(dest, "file_%03d.png")
	cmd := exec.Command("gs", "-sDEVICE=pngalpha", "-o", output, fmt.Sprintf("-r%s", ppi), pdfSrc)
	//cmd.Stdin = os.Stdin
	var errMsg bytes.Buffer
	var outMsg bytes.Buffer
	cmd.Stderr = &errMsg
	cmd.Stdout = &outMsg
	err := cmd.Run()
	if err != nil {
		return errMsg.String(), err
	}
	return outMsg.String(), nil
}

func (g Ghost) ZipDirByPath(dirIn string, zipOut string) error {

	infos, err := ioutil.ReadDir(dirIn)
	if err != nil {
		return errors.Wrapf(err, "")
	}

	var fileIns []string
	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		if !g.isImage(info.Name()) {
			continue
		}
		fileIns = append(fileIns, info.Name())
	}

	f, err := os.Create(zipOut)
	if err != nil {
		return errors.Wrapf(err, "")
	}
	defer f.Close()

	err = g.zipDir(dirIn, fileIns, f)
	if err != nil {
		return errors.Wrapf(err, "")
	}

	return nil
}

func (g Ghost) isImage(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
		return true
	}
	return false
}

func (g Ghost) zipDir(dirIn string, fileIns []string, zipOut io.Writer) error {
	zw := zip.NewWriter(zipOut)
	defer zw.Close()
	for _, file := range fileIns {
		err := g.addFileToZip(zw, dirIn, file)
		if err != nil {
			return errors.Wrapf(err, "")
		}
	}
	return nil
}

func (g Ghost) addFileToZip(zipWriter *zip.Writer, dirIn string, filename string) error {
	//fmt.Printf("===>%s", filepath.Join(dirIn, filename))
	fileToZip, err := os.Open(filepath.Join(dirIn, filename))
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filename
	header.Method = zip.Deflate
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, fileToZip)
	if err != nil {
		return err
	}
	return nil
}

func (g Ghost) ParseOutMsg(msg string) (*GhostScriptInfo, error) {
	token := strings.Split(msg, "\n")
	if token == nil {
		return nil, errors.New("empty")
	}
	info := GhostScriptInfo{}
	size := len(token)
	if size >= 1 {
		info.VersionName = token[0]
	}
	if size >= 4 {
		line := token[3]
		ts := strings.Split(line, " ")
		if len(ts) >= 4 {
			s, err := g.cleanNumber(ts[2])
			if err != nil {
				return nil, err
			}
			e, err := g.cleanNumber(ts[4])
			if err != nil {
				return nil, err
			}
			info.StartFile = s
			info.EndFile = e
		}
	}
	return &info, nil
}

func (g Ghost) cleanNumber(s string) (int, error) {
	tmp := strings.Replace(s, ".", "", -1)
	tmp = strings.Replace(tmp, ",", "", -1)
	tmp = strings.TrimSpace(tmp)
	n, err := strconv.Atoi(tmp)
	if err != nil {
		return 0, err
	}
	return n, nil
}

type GhostScriptInfo struct {
	VersionName string
	StartFile   int
	EndFile     int
}
