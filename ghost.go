package ghost

import (
	"bytes"
	"errors"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Ghost struct {
}

func NewGhost() *Ghost {
	var g Ghost
	return &g
}

func (g Ghost) Convert(pdfSrc, dest string) (string, error) {
	output := filepath.Join(dest, "file_%03d.png")
	cmd := exec.Command("gs", "-sDEVICE=pngalpha", "-o", output, "-r144", pdfSrc)
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
