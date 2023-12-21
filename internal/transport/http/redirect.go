package http

import (
	"context"
	"embed"
	"html/template"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/partyzanex/shortlink/internal/link"
	"github.com/partyzanex/shortlink/pkg/logger"
	"github.com/partyzanex/shortlink/pkg/ptr"
)

const Page = `
<html>
<head>
  <meta charset="UTF-8">
  <title>Redirecting...</title>
</head>
<body>
{{.Counters}}
<script type="text/javascript">
  (function (w, d) { setTimeout(function () { w.location.replace('{{.URL}}') }, 1500) })(window, document);
</script>
</body>
</html>
`

var (
	//go:embed favicon.ico
	Favicon embed.FS

	tpl  *template.Template
	once sync.Once
)

type Params struct {
	Counters string
	URL      string
}

func Execute(w io.Writer, params Params) error {
	var err error

	once.Do(func() {
		tpl, err = template.New("page").Parse(Page)
	})

	if err != nil {
		return errors.Wrap(err, "unable to parse template")
	}

	err = tpl.Execute(w, params)
	if err != nil {
		return errors.Wrap(err, "unable to execute template")
	}

	return nil
}

type LinkService interface {
	Get(ctx context.Context, id *link.ID) (*link.Link, error)
}

func GetRedirectHandler(linkService LinkService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.String(http.StatusBadRequest, "bad request")
			return
		}

		linkInfo, err := linkService.Get(ctx, ptr.Ptr(id))
		if err != nil {
			logger.GetLogger().WithError(err).Error("links.Get(", id, ")")
			ctx.String(http.StatusInternalServerError, "internal error")
			return
		}

		params := Params{
			Counters: "",
			URL:      linkInfo.RawURL(),
		}

		err = Execute(ctx.Writer, params)
		if err != nil {
			logger.GetLogger().WithError(err).Error("Execute ", linkInfo.RawURL())
			ctx.String(http.StatusInternalServerError, "internal error")
			return
		}
	}
}
