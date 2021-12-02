package ali_drive

func (y *Ali) FileUploadPart(reqBody *ApiFileUploadPartRequest) *ApiFileUploadPartResponse {
	url := reqBody.Uri
	req := y.r().
		SetBody(reqBody)

	var resp = new(ApiFileUploadPartResponse)
	if !y.unmarshal(y.processResp(req.Put(url)), &resp) {
		return nil
	}
	return resp
}
