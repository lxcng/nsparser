package parser

type ParserOpts struct {
	Resolution         *int    `json:",omitempty"`
	EpisodeNumberRegex *string `json:",omitempty"`
	MagnetRegex        *string `json:",omitempty"`
	TorrentRegex       *string `json:",omitempty"`
	Translator         *string `json:",omitempty"`
}

func (po *ParserOpts) complete() bool {
	if po.EpisodeNumberRegex == nil {
		return false
	}
	if po.Resolution == nil {
		return false
	}
	if po.MagnetRegex == nil {
		return false
	}
	if po.Translator == nil {
		return false
	}
	return true
}

func (po *ParserOpts) Merge(other *ParserOpts) ParserOpts {
	if po.complete() {
		return *po
	}
	if po.Resolution == nil && other.Resolution != nil {
		po.Resolution = other.Resolution
	}
	if po.EpisodeNumberRegex == nil && other.EpisodeNumberRegex != nil {
		po.EpisodeNumberRegex = other.EpisodeNumberRegex
	}
	if po.MagnetRegex == nil && other.MagnetRegex != nil {
		po.MagnetRegex = other.MagnetRegex
	}
	if po.TorrentRegex == nil && other.TorrentRegex != nil {
		po.TorrentRegex = other.TorrentRegex
	}
	if po.Translator == nil && other.Translator != nil {
		po.Translator = other.Translator
	}
	return *po
}
