package filemerge

import (
	"strings"
)

const (
	defaultMarkerPrefix = "<!-- subikit:"
	defaultMarkerSuffix = " -->"
	defaultClosePrefix  = "<!-- /subikit:"
)

// OpenMarker returns the default opening marker for a section ID.
// Example: sectionID="managed-rules" -> "<!-- subikit:managed-rules -->"
func OpenMarker(sectionID string) string {
	return defaultMarkerPrefix + sectionID + defaultMarkerSuffix
}

// CloseMarker returns the default closing marker for a section ID.
// Example: sectionID="managed-rules" -> "<!-- /subikit:managed-rules -->"
func CloseMarker(sectionID string) string {
	return defaultClosePrefix + sectionID + defaultMarkerSuffix
}

// ExtractMarkdownSection extracts content between default markers for a sectionID.
func ExtractMarkdownSection(content, sectionID string) string {
	return ExtractSectionWithMarkers(content, OpenMarker(sectionID), CloseMarker(sectionID))
}

// ExtractSectionWithMarkers extracts content between custom opening and closing markers.
// Returns empty string if markers are missing or reversed.
func ExtractSectionWithMarkers(content, open, close string) string {
	openIdx := strings.Index(content, open)
	closeIdx := strings.Index(content, close)
	if openIdx == -1 || closeIdx == -1 || closeIdx <= openIdx {
		return ""
	}
	bodyStart := openIdx + len(open)
	return strings.Trim(content[bodyStart:closeIdx], "\r\n")
}

// StripOrphanMarkers removes unpaired opening or closing markers from content.
func StripOrphanMarkers(content, open, close string) string {
	for {
		openIdx := strings.Index(content, open)
		closeIdx := strings.Index(content, close)

		switch {
		case openIdx < 0 && closeIdx < 0:
			return content

		case openIdx < 0 && closeIdx >= 0:
			// Orphan close marker
			content = content[:closeIdx] + content[closeIdx+len(close):]

		case openIdx >= 0 && closeIdx < 0:
			// Orphan open marker
			content = content[:openIdx] + content[openIdx+len(open):]

		case closeIdx < openIdx:
			// Close marker precedes open marker
			content = content[:closeIdx] + content[closeIdx+len(close):]

		default:
			// Valid pair found
			return content
		}
	}
}

// InjectMarkdownSection replaces or appends a marked section in markdown using default subikit markers:
// <!-- subikit:SECTION_ID --> ... <!-- /subikit:SECTION_ID -->
func InjectMarkdownSection(existing, sectionID, content string) string {
	return InjectSectionWithMarkers(existing, OpenMarker(sectionID), CloseMarker(sectionID), content)
}

// InjectSectionWithMarkers replaces or appends a marked section in content using arbitrary markers.
// Rules:
// 1. If content is empty and section exists, the section and its markers are cleanly removed.
// 2. If section already exists, its contents are replaced while preserving everything before and after.
// 3. Any duplicate occurrences of the same section are collapsed.
// 4. If section does not exist, it is appended at the end with clean spacing.
// 5. Orphan markers are stripped beforehand to repair corrupted files.
func InjectSectionWithMarkers(existing, open, close, content string) string {
	// 1. Repair orphan markers
	existing = StripOrphanMarkers(existing, open, close)

	openIdx := strings.Index(existing, open)
	closeIdx := strings.Index(existing, close)

	// 2. Section exists
	if openIdx >= 0 && closeIdx >= 0 && closeIdx > openIdx {
		before := existing[:openIdx]
		after := existing[closeIdx+len(close):]

		// Collapse duplicate blocks if any exist in the remainder
		var preservedAfter strings.Builder
		for {
			dupOpen := strings.Index(after, open)
			if dupOpen < 0 {
				preservedAfter.WriteString(after)
				break
			}
			bodyStart := dupOpen + len(open)
			dupCloseOffset := strings.Index(after[bodyStart:], close)
			if dupCloseOffset < 0 {
				preservedAfter.WriteString(after)
				break
			}
			dupEnd := bodyStart + dupCloseOffset + len(close)
			preservedAfter.WriteString(after[:dupOpen])
			after = after[dupEnd:]
		}
		after = preservedAfter.String()

		// If content is empty, remove the entire section including markers
		if strings.TrimSpace(content) == "" {
			if len(after) > 0 && after[0] == '\n' {
				after = after[1:]
			}
			trimmedBefore := strings.TrimRight(before, "\r\n")
			trimmedAfter := strings.TrimLeft(after, "\r\n")
			if trimmedBefore != "" && trimmedAfter != "" {
				return trimmedBefore + "\n\n" + trimmedAfter
			}
			if trimmedBefore != "" {
				return trimmedBefore + "\n"
			}
			if trimmedAfter != "" {
				return trimmedAfter
			}
			return ""
		}

		var sb strings.Builder
		sb.WriteString(before)
		sb.WriteString(open)
		sb.WriteString("\n")
		sb.WriteString(content)
		if !strings.HasSuffix(content, "\n") {
			sb.WriteString("\n")
		}
		sb.WriteString(close)
		sb.WriteString(after)
		return sb.String()
	}

	// 3. Section not found and content is empty
	if strings.TrimSpace(content) == "" {
		return existing
	}

	// 4. Section not found - append at end
	var sb strings.Builder
	sb.WriteString(existing)
	if existing != "" && !strings.HasSuffix(existing, "\n") {
		sb.WriteString("\n")
	}
	if existing != "" {
		sb.WriteString("\n")
	}
	sb.WriteString(open)
	sb.WriteString("\n")
	sb.WriteString(content)
	if !strings.HasSuffix(content, "\n") {
		sb.WriteString("\n")
	}
	sb.WriteString(close)
	sb.WriteString("\n")
	return sb.String()
}
