package rodutil

import (
	"fmt"
	"strings"
	"sync"

	rod "github.com/Fromsko/rodPro"
	"github.com/Fromsko/rodPro/lib/launcher"
	"github.com/Fromsko/rodPro/lib/proto"
	"github.com/duke-git/lancet/v2/fileutil"
)

// RodClient 定义统一接口
type RodClient interface {
	InitWebClient(wsurl string) *rod.Browser
	SearchParams(page *rod.Page, text string) (*rod.SearchResult, bool)
	SaveHTML(page *rod.Page, path string) bool
	SavePageScreen(pg *rod.Page, path ...string) []byte
	SaveElementScreen(el *rod.Element, path ...string) []byte
	HookResource(page *rod.Page)
	StartWebPage(url string) *rod.Page
}

type option struct {
	Browser *rod.Browser
	mux     *sync.Mutex
	errPut  bool
}

func UseRodTool(options ...func(*option)) RodClient {
	hookOption := &option{
		mux:    &sync.Mutex{},
		errPut: false,
	}

	for _, option := range options {
		option(hookOption)
	}

	return hookOption
}

func (h *option) InitWebClient(wsurl string) (browser *rod.Browser) {
	if wsurl == "" {
		if path, exists := launcher.LookPath(); exists {
			wsurl = launcher.New().Bin(path).MustLaunch()
		}
	}

	if browser = rod.New().ControlURL(wsurl).MustConnect(); browser != nil {
		fmt.Printf("webclient: %s", wsurl)
		h.Browser = browser
	}

	return browser
}

func (h *option) WithErrPut(errPut bool) func(*option) {
	return func(o *option) {
		o.errPut = errPut
	}
}

func (h *option) errPrint(msg string, err error) {
	if h.errPut {
		fmt.Printf("%s err=%s\n", msg, err.Error())
	}
}

func (h *option) StartWebPage(url string) *rod.Page {
	return h.Browser.MustPage(url).MustWaitLoad().MustWaitStable()
}

func (h *option) SearchParams(page *rod.Page, text string) (*rod.SearchResult, bool) {
	search, err := page.Search(text)
	if err != nil {
		return search, false
	}
	return search, search.ResultCount != 0
}

func (h *option) SaveHTML(page *rod.Page, path string) bool {
	h.mux.Lock()
	defer h.mux.Unlock()

	html, err := page.HTML()
	if err != nil {
		h.errPrint("failed to get page HTML", err)
		return false
	}

	if !strings.HasSuffix(path, ".html") {
		path = fmt.Sprintf("%s%d.html", path, GetUnixTime())
	}

	if !fileutil.CreateFile(path) {
		h.errPrint("failed to create file", fmt.Errorf("%s", path))
		return false
	}

	if err := fileutil.WriteStringToFile(path, html, false); err != nil {
		h.errPrint("failed to write HTML to file", err)
		return false
	}
	return true
}

func (h *option) SavePageScreen(pg *rod.Page, path ...string) []byte {
	return pg.MustWaitStable().MustScreenshot(path...)
}

func (h *option) SaveElementScreen(el *rod.Element, path ...string) []byte {
	return el.MustWaitStable().MustScreenshot(path...)
}

func (h *option) HookResource(page *rod.Page) {
	router := page.HijackRequests()

	router.MustAdd("*.png", func(ctx *rod.Hijack) {
		if ctx.Request.Type() == proto.NetworkResourceTypeImage {
			ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
			return
		}
		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	})

	go func() {
		defer func() {
			if r := recover(); r != nil {
				h.errPrint("hijack router panicked", fmt.Errorf("%s", "recover"))
			}
		}()
		router.Run()
	}()
}
