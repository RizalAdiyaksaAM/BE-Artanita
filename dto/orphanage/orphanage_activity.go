package orphanage

type ActivityRequest struct {
	Title          string                  `json:"title" form:"title" validate:"required,min=3"`
	Description    string                  `json:"description" form:"description" validate:"required,min=3"`
	Location       string                  `json:"location" form:"location"`
	Time           string                  `json:"time" form:"time"`
	ActivityImages []ActivityImageRequest `json:"activity_images" form:"activity_images"`
	ActivityVideos []ActivityVideoRequest `json:"activity_videos" form:"activity_videos"`
}

type ActivityImageRequest struct {
	ImageUrl *string `json:"image_url" form:"image_url" validate:"omitempty,required_with=ActivityImages"`
}

type ActivityVideoRequest struct {
	VideoUrl *string `json:"video_url" form:"video_url" validate:"omitempty,required_with=ActivityVideos"`
}


type ActivityResponse struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Location       string `json:"location"`
	Time           string `json:"time"`
	ActivityImages []ActivityImageResponse `json:"activity_images"`
	ActivityVideos []ActivityVideoResponse `json:"activity_videos"`
}

type ActivityImageResponse struct {
	ImageUrl *string `json:"image_url"`
}

type ActivityVideoResponse struct {
	VideoUrl *string `json:"video_url"`
}

type ActivityDocumentResponse struct {
	ActivityImages []ActivityImageResponse `json:"activity_images"`
	ActivityVideos []ActivityVideoResponse `json:"activity_videos"`
}