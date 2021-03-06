package meter

import (
	"testing"

	"github.com/andig/evcc/util/test"
)

func TestMeters(t *testing.T) {
	test.SkipCI(t)

	acceptable := []string{
		"invalid plugin type: ...",
		"missing mqtt broker configuration",
		"mqtt not configured",
		"not a SunSpec device",
		"missing password", // Powerwall
		"connect: no route to host",
		"connect: connection refused",
		"i/o timeout",
	}

	for _, tmpl := range test.ConfigTemplates("meter") {
		t.Run(tmpl.Name, func(t *testing.T) {
			_, err := NewFromConfig(tmpl.Type, tmpl.Config)
			if err != nil && !test.Acceptable(err, acceptable) {
				t.Logf("%s: %+v", tmpl.Name, tmpl.Config)
				t.Error(err)
			}
		})
	}
}
