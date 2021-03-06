package util

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

    <title>{{.}}</title>

    <!-- Everything needs to be kept inline -->
    <style>
      html, body {
        height: 100%;
      }
      body {
        align-items: center;
        background-color: #f7f7f7;
        color: #777;
        display: flex;
        font-family: sans-serif;
        justify-content: center;
        margin: 0;
      }
      a {
        color: #ff6600;
      }
      footer {
        border-top: 2px solid #ddd;
        display: flex;
        padding-top: 12px;
      }
      h1 {
        font-size: 30pt;
        font-weight: normal;
        margin: 0;
      }
      p {
        color: #aaa;
      }
      pre {
        white-space: pre-wrap;
      }
      .content {
        max-width: 800px;
      }
      .logo {
        margin-right: 8px;
        width: 30pt;
      }
      @media (max-width: 839px) {
        .content {
          margin: 20px;
        }
      }
    </style>
  </head>
  <body>
    <div class="content">
      <h1>Something Went Wrong</h1>
      <p>Our server was unable to display the page you requested.</p>
      <p>More technical information may be available below:</p>
      <pre>{{.}}</pre>
      <footer>

        <!-- EVERYTHING needs to be kept inline :P -->
        <img class="logo" src="data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48c3ZnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgd2lkdGg9IjI1NiIgaGVpZ2h0PSIyNTYiIHZpZXdCb3g9IjAgMCA2Ny43MzMgNjcuNzMzIj48ZyBhcmlhLWxhYmVsPSJpNSIgc3R5bGU9ImxpbmUtaGVpZ2h0OjEuMjUiIGZvbnQtd2VpZ2h0PSI0MDAiIGZvbnQtc2l6ZT0iNzQuMDgzIiBmb250LWZhbWlseT0iTG9ic3RlciIgbGV0dGVyLXNwYWNpbmc9IjAiIHdvcmQtc3BhY2luZz0iMCIgc3Ryb2tlLXdpZHRoPSIuMjY1Ij48cGF0aCBkPSJNMjAuMTk5IDIwLjM0NnEtMi40NDUgMC00LjE0OS0xLjcwMy0xLjcwNC0xLjcwNC0xLjcwNC00LjE1IDAtMi40NDQgMS43MDQtNC4xNDggMS43MDQtMS43NzggNC4xNDktMS43NzggMi40NDUgMCA0LjE0OSAxLjc3OCAxLjc3OCAxLjcwNCAxLjc3OCA0LjE0OSAwIDIuNDQ1LTEuNzc4IDQuMTQ5LTEuNzA0IDEuNzAzLTQuMTUgMS43MDN6TTE1LjE2IDYxLjgzM3EtMy41NTYgMC01Ljc3OC0yLjIyMi0yLjE0OS0yLjIyMy0yLjE0OS02LjY2OCAwLTEuODUyLjU5My00LjgxNWw1LjAzOC0yMy43ODFoMTAuNjY4bC01LjMzNCAyNS4xODhxLS4yOTcgMS4xMTItLjI5NyAyLjM3MSAwIDEuNDgyLjY2NyAyLjE0OC43NC41OTMgMi4zNy41OTMgMS4zMzQgMCAyLjM3MS0uNDQ0LS4yOTYgMy43MDQtMi42NjcgNS43MDQtMi4yOTYgMS45MjYtNS40ODIgMS45MjZ6IiBzdHlsZT0iLWlua3NjYXBlLWZvbnQtc3BlY2lmaWNhdGlvbjpMb2JzdGVyIiBmaWxsPSIjZjYwIi8+PHBhdGggZD0iTTM3LjkwNCA2MS43NTlxLTYuNTIgMC0xMC4wMDItMi44MTUtMy40ODItMi44MTUtMy40ODItOC4wMDEgMC00LjY2NyAyLjY2Ny03LjQ4MyAyLjY2Ny0yLjg4OSA2Ljk2NC0yLjg4OSA0LjUyIDAgNS4yNiAzLjc3OC0zLjAzNy4zNy00Ljg5IDIuMjIzLTEuNzc3IDEuNzc4LTEuNzc3IDQuNjY3IDAgMi41OTMgMS41NTUgNC4wNzUgMS41NTYgMS40ODEgNC4xNSAxLjQ4MSA0LjI5NiAwIDcuMDM3LTMuODUyIDIuNzQxLTMuODUyIDIuNzQxLTkuODUzIDAtNS4xMTItMi4xNDgtNy44NTMtMi4xNDktMi44MTUtNi4wMDEtMi44MTUtMy45MjYgMC05LjcwNSAzLjMzNGwtLjk2My0uNjY3IDUuNzA0LTI4LjA3OHE4LjM3Mi45NjQgMTIuNDQ2Ljk2NCA2Ljc0MiAwIDEyLjg5LTIuMDc1LjE0OSAxLjYzLjE0OSAyLjQ0NSAwIDMuODUyLTIuMDc0IDYuMjIzLTIgMi4zNy02LjQ0NiAyLjM3LTEuNzAzIDAtNC4wNzQtLjM3LTIuMzctLjM3LTQuNTItLjgxNS0xLjg1MS0uMzctMy43NzctLjY2NmwtMi45NjQgMTMuNDA5cTQuODE2LTIuMjk3IDkuMTg3LTIuMjk3IDYuMDc0IDAgOS41NTYgNC4xNDkgMy41NTYgNC4wNzQgMy41NTYgMTEuMTEyIDAgNi4wMDEtMi42NjcgMTAuNTk0LTIuNTkzIDQuNTkzLTcuNDA4IDcuMTg2LTQuNzQxIDIuNTE5LTEwLjk2NCAyLjUxOXoiIHN0eWxlPSItaW5rc2NhcGUtZm9udC1zcGVjaWZpY2F0aW9uOkxvYnN0ZXIiIGZpbGw9InB1cnBsZSIvPjwvZz48L3N2Zz4NCg==">

        <!-- Branding -->
        <p>
          Powered by:<br>
          <a href="https://github.com/nathan-osman/i5">i5 (reverse proxy)</a>
        </p>
      </footer>
    </div>
  </body>
</html>
`))

// RenderError renders the error template with the specified error.
func RenderError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	buff := &bytes.Buffer{}
	errorTemplate.Execute(buff, message)
	b := buff.Bytes()
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(b)
}
