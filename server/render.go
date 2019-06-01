package server

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"
)

var errorTemplate = template.Must(template.New("").Parse(`<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.StatusCode}} {{.StatusText}}</title>
    <style>
      body { background: #eee; color: #777; font-family: sans-serif; margin: 2em; }
      footer { border-top: 2px solid #ddd; color: #aaa; display: flex; margin-top: 4em; padding-top: 1em; }
      h1 { font-size: 3em; margin-top: 0; }
      h1 .subtitle { color: #aaa; display: block; font-size: 75%; font-weight: normal; }
      p { font-size: 1em; }
      .logo { align-self: center; margin-right: 1em; max-height: 2em; }
	</style>
  </head>
  <body>
	<h1>
      An Error Has Occurred
      <span class="subtitle">{{.StatusCode}} {{.StatusText}}</span>
    </h1>
    <p>Unfortunately, something went wrong while processing your request.</p>
    <p>Please try again later.</p>
    <footer>
      <img class="logo" src="data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48c3ZnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgd2lkdGg9IjI1NiIgaGVpZ2h0PSIyNTYiIHZpZXdCb3g9IjAgMCA2Ny43MzMgNjcuNzMzIj48ZyBhcmlhLWxhYmVsPSJpNSIgc3R5bGU9ImxpbmUtaGVpZ2h0OjEuMjUiIGZvbnQtd2VpZ2h0PSI0MDAiIGZvbnQtc2l6ZT0iNzQuMDgzIiBmb250LWZhbWlseT0iTG9ic3RlciIgbGV0dGVyLXNwYWNpbmc9IjAiIHdvcmQtc3BhY2luZz0iMCIgc3Ryb2tlLXdpZHRoPSIuMjY1Ij48cGF0aCBkPSJNMjAuMTk5IDIwLjM0NnEtMi40NDUgMC00LjE0OS0xLjcwMy0xLjcwNC0xLjcwNC0xLjcwNC00LjE1IDAtMi40NDQgMS43MDQtNC4xNDggMS43MDQtMS43NzggNC4xNDktMS43NzggMi40NDUgMCA0LjE0OSAxLjc3OCAxLjc3OCAxLjcwNCAxLjc3OCA0LjE0OSAwIDIuNDQ1LTEuNzc4IDQuMTQ5LTEuNzA0IDEuNzAzLTQuMTUgMS43MDN6TTE1LjE2IDYxLjgzM3EtMy41NTYgMC01Ljc3OC0yLjIyMi0yLjE0OS0yLjIyMy0yLjE0OS02LjY2OCAwLTEuODUyLjU5My00LjgxNWw1LjAzOC0yMy43ODFoMTAuNjY4bC01LjMzNCAyNS4xODhxLS4yOTcgMS4xMTItLjI5NyAyLjM3MSAwIDEuNDgyLjY2NyAyLjE0OC43NC41OTMgMi4zNy41OTMgMS4zMzQgMCAyLjM3MS0uNDQ0LS4yOTYgMy43MDQtMi42NjcgNS43MDQtMi4yOTYgMS45MjYtNS40ODIgMS45MjZ6IiBzdHlsZT0iLWlua3NjYXBlLWZvbnQtc3BlY2lmaWNhdGlvbjpMb2JzdGVyIiBmaWxsPSIjZjYwIi8+PHBhdGggZD0iTTM3LjkwNCA2MS43NTlxLTYuNTIgMC0xMC4wMDItMi44MTUtMy40ODItMi44MTUtMy40ODItOC4wMDEgMC00LjY2NyAyLjY2Ny03LjQ4MyAyLjY2Ny0yLjg4OSA2Ljk2NC0yLjg4OSA0LjUyIDAgNS4yNiAzLjc3OC0zLjAzNy4zNy00Ljg5IDIuMjIzLTEuNzc3IDEuNzc4LTEuNzc3IDQuNjY3IDAgMi41OTMgMS41NTUgNC4wNzUgMS41NTYgMS40ODEgNC4xNSAxLjQ4MSA0LjI5NiAwIDcuMDM3LTMuODUyIDIuNzQxLTMuODUyIDIuNzQxLTkuODUzIDAtNS4xMTItMi4xNDgtNy44NTMtMi4xNDktMi44MTUtNi4wMDEtMi44MTUtMy45MjYgMC05LjcwNSAzLjMzNGwtLjk2My0uNjY3IDUuNzA0LTI4LjA3OHE4LjM3Mi45NjQgMTIuNDQ2Ljk2NCA2Ljc0MiAwIDEyLjg5LTIuMDc1LjE0OSAxLjYzLjE0OSAyLjQ0NSAwIDMuODUyLTIuMDc0IDYuMjIzLTIgMi4zNy02LjQ0NiAyLjM3LTEuNzAzIDAtNC4wNzQtLjM3LTIuMzctLjM3LTQuNTItLjgxNS0xLjg1MS0uMzctMy43NzctLjY2NmwtMi45NjQgMTMuNDA5cTQuODE2LTIuMjk3IDkuMTg3LTIuMjk3IDYuMDc0IDAgOS41NTYgNC4xNDkgMy41NTYgNC4wNzQgMy41NTYgMTEuMTEyIDAgNi4wMDEtMi42NjcgMTAuNTk0LTIuNTkzIDQuNTkzLTcuNDA4IDcuMTg2LTQuNzQxIDIuNTE5LTEwLjk2NCAyLjUxOXoiIHN0eWxlPSItaW5rc2NhcGUtZm9udC1zcGVjaWZpY2F0aW9uOkxvYnN0ZXIiIGZpbGw9InB1cnBsZSIvPjwvZz48L3N2Zz4NCg==">
      <div>
        Powered by:<br>
        <a href="https://github.com/nathan-osman/i5">i5 (reverse proxy)</a>
        &mdash; Copyright 2019 by Nathan Osman
      </div>
    </footer>
  </body>
</html>
`))

func renderErrorTemplate(w http.ResponseWriter, r *http.Request, statusCode int) {
	buff := &bytes.Buffer{}
	errorTemplate.Execute(buff,
		struct {
			StatusCode int
			StatusText string
		}{
			StatusCode: statusCode,
			StatusText: http.StatusText(statusCode),
		},
	)
	b := buff.Bytes()
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(b)
}
