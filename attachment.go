package nerimity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// cdnBaseURL is the Nerimity CDN base. The CDN is a separate service from the
// main API and is not affected by APIURLOverride.
const cdnBaseURL = "https://cdn.nerimity.com"

// AttachmentBuilder uploads a file to the Nerimity CDN and produces a file ID
// that can be attached to a message via MessageOptions.NerimityCdnFileID.
type AttachmentBuilder struct {
	reader io.Reader
	name   string
}

// NewAttachment creates an AttachmentBuilder from an io.Reader. name is the
// filename to present to the CDN.
func NewAttachment(r io.Reader, name string) *AttachmentBuilder {
	return &AttachmentBuilder{reader: r, name: name}
}

// NewAttachmentFromFile creates an AttachmentBuilder from a file on disk. The
// file is read when Build is called.
func NewAttachmentFromFile(path string) (*AttachmentBuilder, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("nerimity: opening attachment: %w", err)
	}
	return &AttachmentBuilder{reader: f, name: fileBase(path)}, nil
}

func fileBase(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[i+1:]
		}
	}
	return path
}

// Build uploads the file, associates it with the given channel, and returns the
// CDN file ID. Pass that ID as MessageOptions.NerimityCdnFileID when sending a
// message on the same channel. If the reader is an io.Closer it is closed after
// the upload.
func (a *AttachmentBuilder) Build(ctx context.Context, client *Client, channel *Channel) (string, error) {
	if c, ok := a.reader.(io.Closer); ok {
		defer c.Close()
	}

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, err := mw.CreateFormFile("file", a.name)
	if err != nil {
		return "", fmt.Errorf("nerimity: building upload form: %w", err)
	}
	if _, err := io.Copy(fw, a.reader); err != nil {
		return "", fmt.Errorf("nerimity: reading attachment: %w", err)
	}
	if err := mw.Close(); err != nil {
		return "", fmt.Errorf("nerimity: finalising upload form: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cdnBaseURL+"/upload", &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("nerimity: uploading attachment: %w", err)
	}
	uploadBody, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", &APIError{Status: resp.StatusCode, Body: string(uploadBody)}
	}
	var upload rawCDNUpload
	if err := json.Unmarshal(uploadBody, &upload); err != nil {
		return "", fmt.Errorf("nerimity: decoding upload response: %w", err)
	}

	// Associate the uploaded file with the channel.
	saveURL := fmt.Sprintf("%s/attachments/%s/%s", cdnBaseURL, channel.ID, upload.FileID)
	if err := client.doJSON(ctx, http.MethodPost, saveURL, upload, nil, ""); err != nil {
		return "", fmt.Errorf("nerimity: saving attachment: %w", err)
	}
	return upload.FileID, nil
}
