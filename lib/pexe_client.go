package lib

const CUSTOM string = "CUSTOM"
const AWS string = "AWS"

type PexecClient struct {

	TargetServers []Server

	Provider string

	Parallel bool

	Batch bool

	BatchSize int

}

func (client *PexecClient) getDefaults()  {

	if client.Provider == "" {
		client.Provider = CUSTOM
	}
}
