package pipeline

func New() Pipeline {
	return &pipelineImpl{
		scraper:          nil,
		assetDownloader:  nil,
		assetURLFetcher:  nil,
		jobIndexer:       nil,
		runContext:       nil,
		cancelRunContext: nil,
	}
}
