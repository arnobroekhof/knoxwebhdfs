package knoxwebhdfs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

var (
	ErrorNotFound         = errors.New("file not found")
	ErrorIsDirectory      = errors.New("is directory")
	ErrorMoreThenOneEntry = errors.New("more then one entry found")
	ErrorExists           = errors.New("path or file already exists")
)

func (c *Client) Get(path string) (io.Reader, error) {
	stat, err := c.List(path)
	if err != nil {
		return nil, fmt.Errorf("%v %w", ErrorNotFound, err)
	}

	if len(stat.FileStatus) > 1 {
		return nil, ErrorMoreThenOneEntry
	}

	if stat.FileStatus[0].Type == "DIRECTORY" {
		return nil, ErrorIsDirectory
	}

	reqPath := fmt.Sprintf("%s/%s?op=OPEN", c.url.String(), path)
	req, err := http.NewRequest("GET", reqPath, nil)
	if err != nil {
		return nil, err
	}

	if c.conf.AuthType == AuthTypeBasic {
		req.SetBasicAuth(c.conf.Username, c.conf.Password)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp.Body, nil

}

func (c *Client) mkdir(path string) error {
	exists, err := c.dirExists(path)
	if err != nil {
		return err
	}
	if !exists {
		reqPath := fmt.Sprintf("%s/%s?op=MKDIRS", c.url.String(), path)
		req, err := http.NewRequest("PUT", reqPath, nil)
		if err != nil {
			return err
		}

		if c.conf.AuthType == AuthTypeBasic {
			req.SetBasicAuth(c.conf.Username, c.conf.Password)
		}

		resp, err := c.client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}
	return nil
}

func (c *Client) Put(path string, reader io.Reader) error {

	baseDir := filepath.Base(path)
	err := c.mkdir(baseDir)
	if err != nil {
		return err
	}

	reqPath := fmt.Sprintf("%s/%s?op=CREATE&overwrite=true", c.url.String(), path)
	req, err := http.NewRequest("PUT", reqPath, reader)
	if err != nil {
		return err
	}

	if c.conf.AuthType == AuthTypeBasic {
		req.SetBasicAuth(c.conf.Username, c.conf.Password)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (c *Client) dirExists(path string) (bool, error) {
	// check the directory if exists
	reqPath := fmt.Sprintf("%s/%s?op=GETFILESTATUS", c.url.String(), path)
	req, err := http.NewRequest("GET", reqPath, nil)
	if err != nil {
		return false, err
	}

	if c.conf.AuthType == AuthTypeBasic {
		req.SetBasicAuth(c.conf.Username, c.conf.Password)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}

	var fileStatus DirFileStatus
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(respBody, &fileStatus)
	if err != nil {
		return false, err
	}

	if fileStatus.FileStatus.Type == "DIRECTORY" {
		return true, nil
	}

	return false, nil
}

func (c *Client) List(path string) (list FileStatuses, err error) {
	reqPath := fmt.Sprintf("%s/%s?op=LISTSTATUS", c.url.String(), path)
	req, err := http.NewRequest("GET", reqPath, nil)
	if err != nil {
		return list, err
	}

	if c.conf.AuthType == AuthTypeBasic {
		req.SetBasicAuth(c.conf.Username, c.conf.Password)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return list, err
	}

	defer resp.Body.Close()

	var fileStatuses WebHdfsFileStatuses
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return list, err
	}

	err = json.Unmarshal(respBody, &fileStatuses)
	if err != nil {
		return list, err
	}
	return fileStatuses.FileStatuses, err
}
