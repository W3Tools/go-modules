package gmsui

// type SuiSignatureCombineClient struct {
// 	ClientUrl    string
// 	PublicKeyMap []SuiPubkeyWeightPair
// 	Threshold    uint16
// }

// type SuiCombineSignatureRequest struct {
// 	Threshold      uint16                              `json:"threshold"`
// 	Signatures     []string                            `json:"signatures"`
// 	PublicKeyPairs []SuiCombineSignatureRequestKeyPair `json:"publicKeyPairs"`
// }
// type SuiCombineSignatureRequestKeyPair struct {
// 	B64PublicKey string `json:"b64PublicKey"`
// 	Weight       uint8  `json:"weight"`
// }

// type SuiCombineSignatureResponse struct {
// 	Code int64                           `json:"code"`
// 	Msg  string                          `json:"msg"`
// 	Data SuiCombineSignatureResponseData `json:"data"`
// }

// type SuiCombineSignatureResponseData struct {
// 	Address    string `json:"address"`
// 	Serialized string `json:"serialized"`
// }

// func NewSuiSignatureCombineClient(publicKeyMap []SuiPubkeyWeightPair, threshold uint16) *SuiSignatureCombineClient {
// 	return &SuiSignatureCombineClient{
// 		ClientUrl:    "https://api-tools-sui.w3tools.app/v1/multisig/signature/combine",
// 		PublicKeyMap: publicKeyMap,
// 		Threshold:    threshold,
// 	}
// }

// func (c *SuiSignatureCombineClient) NewCombineSignatureRequestData(signatures []string) *SuiCombineSignatureRequest {
// 	ret := &SuiCombineSignatureRequest{
// 		Threshold:  c.Threshold,
// 		Signatures: signatures,
// 	}

// 	for _, k := range c.PublicKeyMap {
// 		ret.PublicKeyPairs = append(ret.PublicKeyPairs, SuiCombineSignatureRequestKeyPair{
// 			B64PublicKey: k.PublicKey,
// 			Weight:       k.Weight,
// 		})
// 	}

// 	return ret
// }

// func (c *SuiSignatureCombineClient) TryGetCombineSignatures(signatures []string) (*SuiCombineSignatureResponseData, error) {
// 	dataStruct := c.NewCombineSignatureRequestData(signatures)
// 	data, err := json.Marshal(dataStruct)
// 	if err != nil {
// 		return nil, fmt.Errorf("json.Marshal %v", err)
// 	}

// 	req, err := http.NewRequest("POST", c.ClientUrl, bytes.NewBuffer(data))
// 	if err != nil {
// 		return nil, fmt.Errorf("http.NewRequest %v", err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	httpClient := &http.Client{}
// 	result, err := httpClient.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("httpClient.Do %v", err)
// 	}
// 	defer result.Body.Close()

// 	if result.StatusCode != 201 {
// 		return nil, fmt.Errorf("tryGetCombineSignatures http code: %v", result.StatusCode)
// 	}

// 	body, _ := io.ReadAll(result.Body)
// 	ret := &SuiCombineSignatureResponse{}
// 	err = json.Unmarshal(body, ret)
// 	if err != nil {
// 		return nil, fmt.Errorf("json.Unmarshal %v", err)
// 	}
// 	return &ret.Data, nil
// }
