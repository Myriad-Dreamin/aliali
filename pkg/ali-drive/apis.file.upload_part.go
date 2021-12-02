package ali_drive

func (y *Ali) FileUploadPart(reqBody *ApiFileUploadPartRequest) *ApiFileUploadPartResponse {
	url := reqBody.Uri
	req := y.r(y.uploadClient).
		SetBody(reqBody.Reader)

	var resp = new(ApiFileUploadPartResponse)
	if y.processResp(req.Put(url)) == nil {
		return nil
	}
	return resp
}
