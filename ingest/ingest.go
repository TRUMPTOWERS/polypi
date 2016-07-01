package ingest

// Ingest is what the ingested data format looks like
type Ingest struct {
	Pies []struct {
		ID            int64    `json:"id"`
		Name          string   `json:"name"`
		ImageURL      string   `json:"image_url"`
		PricePerSlice float64  `json:"price_per_slice"`
		Slices        int      `json:"slices"`
		Labels        []string `json:"labels"`
	} `json:"pies"`
}
