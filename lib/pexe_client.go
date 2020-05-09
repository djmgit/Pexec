package lib


type PexecClient struct {

	TargetServers []Server

	Provider string

	Parallel bool

	Batch bool

	BatchSize int

}
