package utils

import (
    "encoding/base64"
    "github.com/jaytaylor/html2text"
    "strings"
)

// FormatHTML2Text converts HTML to plain text
func FormatHTML2Text(html string) (string, error) {
    html = strings.Replace(html, "\r\n", "<br />", -1)
    html = strings.Replace(html, "\n", "<br />", -1)
    return html2text.FromString(html, html2text.Options{ PrettyTables: true })
}

// FormatHTML2Base64 converts HTML to base64-decoded plain text
func FormatHTML2Base64(html string) (string, error) {
    if text, err := FormatHTML2Text(html); err == nil {
        return base64.StdEncoding.EncodeToString([]byte(text)), nil
    }

    return "", nil
}
