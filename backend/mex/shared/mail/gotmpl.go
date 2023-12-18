package mail

import (
	"bytes"
	tplHtml "html/template"
	tplText "text/template"
)

func InterpolateGoPlain(template string, data any) (string, error) {
	t, err := tplText.New("email").Parse(template)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func InterpolateGoHTML(template string, data any) (string, error) {
	t, err := tplHtml.New("email").Parse(template)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
