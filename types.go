package knoxwebhdfs

type WebHdfsFileStatuses struct {
	FileStatuses FileStatuses `json:"FileStatuses"`
}

type FileStatuses struct {
	FileStatus []FileStatus `json:"FileStatus"`
}

type DirFileStatus struct {
	FileStatus FileStatus `json:"FileStatus"`
}

type FileStatus struct {
	AccessTime       int64  `json:"accessTime"`
	BlockSize        int    `json:"blockSize"`
	Group            string `json:"group"`
	Length           int    `json:"length"`
	ModificationTime int64  `json:"modificationTime"`
	Owner            string `json:"owner"`
	PathSuffix       string `json:"pathSuffix"`
	Permission       string `json:"permission"`
	Replication      int    `json:"replication"`
	Type             string `json:"type"`
}
