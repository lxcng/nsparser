package adapter

type ClientOpts struct {
	Os             *string `json:",omitempty"`
	TorrentClient  *string `json:",omitempty"`
	DownloadFolder *string `json:",omitempty"`
}

func (co *ClientOpts) complete() bool {
	if co.Os == nil {
		return false
	}
	if co.TorrentClient == nil {
		return false
	}
	if co.DownloadFolder == nil {
		return false
	}
	return true
}

func (co *ClientOpts) compile() interface{} {

	return nil
}

func (co *ClientOpts) Merge(other *ClientOpts) ClientOpts {
	if co.complete() {
		return *co
	}
	if co.Os == nil && other.Os != nil {
		co.Os = other.Os
	}
	if co.TorrentClient == nil && other.TorrentClient != nil {
		co.TorrentClient = other.TorrentClient
	}
	if co.DownloadFolder == nil && other.DownloadFolder != nil {
		co.DownloadFolder = other.DownloadFolder
	}
	return *co
}
