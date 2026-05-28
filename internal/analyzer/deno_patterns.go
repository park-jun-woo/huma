//ff:func feature=analyzer type=helper control=sequence
//ff:what Provides compiled regex patterns for detecting response status codes in Deno/Supabase Edge Functions
package analyzer

import "regexp"

var denoPatterns = []*regexp.Regexp{
	regexp.MustCompile(`new Response\(.*\{\s*status:\s*(\d+)`),
	regexp.MustCompile(`Response\.json\(.*\{\s*status:\s*(\d+)`),
	regexp.MustCompile(`Response\.redirect\([^,]+,\s*(\d+)`),
}

var denoImplicitJson200 = regexp.MustCompile(`Response\.json\(`)

var denoImplicitResponse200 = regexp.MustCompile(`new Response\([^,]*\)`)

var denoImplicitRedirect302 = regexp.MustCompile(`Response\.redirect\(\s*['"\x60]`)
