package translator

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Translator 翻译接口
type Translator interface {
	Translate(text string) (string, error)
	TranslateWithResult(text string) (*TranslationResult, error)
	Extract(result string) (*TranslationResult, error)
}

// Option 配置选项函数类型
type Option func(*translator)

// translator 实现 Translator 接口
type translator struct {
	client    *http.Client
	host      string
	uri       string
	appid     string
	Secret    string
	apiKey    string
	fromLang  string
	toLang    string
	httpProto string
}

// Config 配置结构体（用于结构体方式初始化）
type Config struct {
	Host      string
	URI       string
	AppID     string
	Secret    string
	APIKey    string
	FromLang  string
	ToLang    string
	HTTPProto string
}

// TranslationResult 翻译结果结构体
type TranslationResult struct {
	Source string `json:"src"` // 源文本
	Target string `json:"dst"` // 翻译结果
}

// New 创建翻译实例（支持 Option 方式）
func New(options ...Option) Translator {
	t := &translator{
		client:    &http.Client{},
		host:      "ntrans.xfyun.cn",
		uri:       "/v2/ots",
		appid:     "**",
		Secret:    "****",
		apiKey:    "****",
		fromLang:  "cn",
		toLang:    "en",
		httpProto: "HTTP/1.1",
	}

	for _, opt := range options {
		opt(t)
	}

	return t
}

// NewWithConfig 创建翻译实例（支持结构体方式）
func NewWithConfig(config Config) Translator {
	return &translator{
		client:    &http.Client{},
		host:      config.Host,
		uri:       config.URI,
		appid:     config.AppID,
		Secret:    config.Secret,
		apiKey:    config.APIKey,
		fromLang:  config.FromLang,
		toLang:    config.ToLang,
		httpProto: config.HTTPProto,
	}
}

// WithAppID 设置 AppID
func WithAppID(appid string) Option {
	return func(t *translator) {
		t.appid = appid
	}
}

// WithSecret 设置 Secret
func WithSecret(secret string) Option {
	return func(t *translator) {
		t.Secret = secret
	}
}

// WithAPIKey 设置 API Key
func WithAPIKey(apiKey string) Option {
	return func(t *translator) {
		t.apiKey = apiKey
	}
}

// WithFromLang 设置源语言
func WithFromLang(lang string) Option {
	return func(t *translator) {
		t.fromLang = lang
	}
}

// WithToLang 设置目标语言
func WithToLang(lang string) Option {
	return func(t *translator) {
		t.toLang = lang
	}
}

// WithHost 设置主机地址
func WithHost(host string) Option {
	return func(t *translator) {
		t.host = host
	}
}

// WithURI 设置请求路径
func WithURI(uri string) Option {
	return func(t *translator) {
		t.uri = uri
	}
}

// WithHTTPProto 设置 HTTP 协议版本
func WithHTTPProto(proto string) Option {
	return func(t *translator) {
		t.httpProto = proto
	}
}

// Translate 翻译文本
func (t *translator) Translate(text string) (string, error) {
	data := []byte(text)
	param := map[string]any{
		"common": map[string]any{
			"app_id": t.appid,
		},
		"business": map[string]any{
			"from": t.fromLang,
			"to":   t.toLang,
		},
		"data": map[string]any{
			"text": base64.StdEncoding.EncodeToString(data),
		},
	}

	jsonData, err := json.Marshal(param)
	if err != nil {
		return "", fmt.Errorf("marshal json failed: %w", err)
	}

	requestBody := bytes.NewBuffer(jsonData)

	currentTime := time.Now().UTC().Format(time.RFC1123)
	digest := fmt.Sprintf("SHA-256=%s", signBody(string(jsonData)))

	signature := generateSignature(
		t.host,
		currentTime,
		"POST",
		t.uri,
		t.httpProto,
		digest,
		t.Secret,
	)

	authHeader := fmt.Sprintf(`api_key="%s", algorithm="%s", headers="host date request-line digest", signature="%s"`,
		t.apiKey, "hmac-sha256", signature)

	request, err := http.NewRequest("POST", "http://"+t.host+t.uri, requestBody)
	if err != nil {
		return "", fmt.Errorf("create request failed: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Host", t.host)
	request.Header.Set("Accept", "application/json,version=1.0")
	request.Header.Set("Date", currentTime)
	request.Header.Set("Digest", digest)
	request.Header.Set("Authorization", authHeader)

	response, err := t.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("read response failed: %w", err)
	}

	return string(body), nil
}

// Extract 从API响应中提取翻译结果
func (t *translator) Extract(result string) (*TranslationResult, error) {
	// 定义完整的响应结构体
	type FullResponse struct {
		Data struct {
			Result struct {
				TransResult struct {
					Src string `json:"src"`
					Dst string `json:"dst"`
				} `json:"trans_result"`
			} `json:"result"`
		} `json:"data"`
	}

	var fullResponse FullResponse
	if err := json.Unmarshal([]byte(result), &fullResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 检查是否有翻译结果
	if fullResponse.Data.Result.TransResult.Dst == "" {
		return nil, fmt.Errorf("no translation result found")
	}

	transResult := fullResponse.Data.Result.TransResult
	return &TranslationResult{
		Source: transResult.Src,
		Target: transResult.Dst,
	}, nil
}

// TranslateWithResult 翻译并提取字段
func (t *translator) TranslateWithResult(text string) (*TranslationResult, error) {
	data, err := t.Translate(text)
	if err != nil {
		return nil, err
	}
	return t.Extract(data)
}

// generateSignature 生成签名
func generateSignature(host, date, httpMethod, requestUri, httpProto, digest, secret string) string {
	var signatureStr string
	if len(host) != 0 {
		signatureStr = "host: " + host + "\n"
	}
	signatureStr += "date: " + date + "\n"
	signatureStr += httpMethod + " " + requestUri + " " + httpProto + "\n"
	signatureStr += "digest: " + digest
	return hmacsign(signatureStr, secret)
}

// hmacsign 计算 HMAC-SHA256 并 Base64 编码
func hmacsign(data, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

// signBody 计算请求体的 SHA-256 并 Base64 编码
func signBody(data string) string {
	sha := sha256.New()
	sha.Write([]byte(data))
	encodeData := sha.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}
