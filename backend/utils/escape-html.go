package utils

import "github.com/microcosm-cc/bluemonday"

var policy *bluemonday.Policy

func InitBluemonday() {
	policy = bluemonday.UGCPolicy()
	policy.AllowDataAttributes()

	// Text formatting
	policy.AllowElements("b", "i", "u", "strong", "em", "small", "mark")
	policy.AllowElements("br", "hr")

	// Lists
	policy.AllowElements("ul", "ol", "li", "dl", "dt", "dd")

	// Tables
	policy.AllowElements("table", "thead", "tbody", "tfoot", "tr", "th", "td")

	// Code / preformatted
	policy.AllowElements("code", "pre")

	// Blockquote
	policy.AllowElements("blockquote")

	// Forms
	policy.AllowElements("input", "textarea", "button", "label", "select", "option", "fieldset", "legend")

	// Multimedia
	policy.AllowElements("audio", "video")

	// Custom elements for frontend rendering
	policy.AllowElements("tag")

	// SVG elements for math/graphics
	policy.AllowElements("svg", "path", "line", "polygon", "polyline", "rect", "circle", "ellipse", "text", "g", "style")

	// Global attributes
	policy.AllowAttrs(
		"class", "id", "style", "tabindex", "aria-hidden", "display",
		"xmlns", "encoding", "accent", "width", "height", "type", "title",
	).Globally()

	// Custom color attributes
	policy.AllowAttrs(
		"grey", "blue", "red", "green", "yellow", "purple", "orange", "teal", "pink", "primary",
	).Globally()

	// SVG-specific attributes
	policy.AllowAttrs(
		"viewBox", "fill-rule", "d", "x", "y", "cx", "cy", "r", "rx", "ry", "points", "stroke", "stroke-width", "fill", "preserveAspectRatio", "xmlns",
	).OnElements("svg", "path", "line", "polygon", "polyline", "rect", "circle", "ellipse", "text")

	// Links
	policy.AllowAttrs("href", "rel").OnElements("a")

	// Form attributes
	policy.AllowAttrs("checked").OnElements("input")

	// Disable automatic rel="nofollow"
	policy.RequireNoFollowOnLinks(false)
}

// EscapeHTML sanitizes user input to prevent XSS attacks
func EscapeHTML(input *string) string {
	if input == nil || *input == "" {
		return ""
	}
	return policy.Sanitize(*input)
}
