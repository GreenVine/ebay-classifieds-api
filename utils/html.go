package utils

import (
    "encoding/base64"
    "jaytaylor.com/html2text"
    "strings"
)

func FormatHtml2Text(html string) (string, error) {
    html = strings.Replace(html, "\r\n", "<br />", -1)
    html = strings.Replace(html, "\n", "<br />", -1)
    return html2text.FromString(html, html2text.Options{ PrettyTables: true })
}

func FormatHtml2Base64(html string) (string, error) {
    if text, err := FormatHtml2Text(html); err == nil {
        return base64.StdEncoding.EncodeToString([]byte(text)), nil
    } else {
        return "", nil
    }
}
