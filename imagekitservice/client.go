package imagekitservice

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	imagekitsdk "github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

const (
	UploadURL            = "https://upload.imagekit.io/api/v2/files/upload"
	UploadAuthTTL        = 5 * time.Minute
	MaxFileSizeBytes     = int64(5 * 1024 * 1024)
	MaxImageDimension    = 4096
	profileImageRoot     = "/profile-images"
	maxOriginalNameBytes = 255
)

var UploadChecks = fmt.Sprintf(
	"'file.mime' IN ['image/jpeg', 'image/png', 'image/webp'] AND 'file.size' <= %d AND 'mediaMetadata.width' <= %d AND 'mediaMetadata.height' <= %d",
	MaxFileSizeBytes,
	MaxImageDimension,
	MaxImageDimension,
)

var (
	ErrInvalidUID              = errors.New("invalid Firebase UID")
	ErrInvalidFilename         = errors.New("invalid image filename")
	ErrInvalidProfileImagePath = errors.New("invalid profile image path")
	ErrUnsupportedExtension    = errors.New("unsupported image extension")
	ErrAssetVerificationFailed = errors.New("profile image asset verification failed")
)

// Config contains the server-only ImageKit credentials and public CDN endpoint.
type Config struct {
	PublicKey   string
	PrivateKey  string
	URLEndpoint string
}

// UploadPayload is the complete set of non-file multipart fields signed into a
// V2 client-upload JWT. Values are strings because ImageKit requires upload
// parameters to be stringified in V2 JWT claims and multipart form fields.
type UploadPayload struct {
	FileName          string `json:"fileName"`
	Folder            string `json:"folder"`
	UseUniqueFileName string `json:"useUniqueFileName"`
	OverwriteFile     string `json:"overwriteFile"`
	IsPrivateFile     string `json:"isPrivateFile"`
	IsPublished       string `json:"isPublished"`
	Checks            string `json:"checks"`
}

// UploadAuth contains the short-lived credentials and immutable fields needed
// for a direct browser upload to ImageKit Upload API V2.
type UploadAuth struct {
	Token         string        `json:"token"`
	ExpiresAt     int64         `json:"expiresAt"`
	UploadURL     string        `json:"uploadUrl"`
	UploadPayload UploadPayload `json:"uploadPayload"`
}

type VerifiedProfileImageAsset struct {
	FileID   string
	FilePath string
}

type Client struct {
	sdk         *imagekitsdk.Client
	publicKey   string
	privateKey  string
	urlEndpoint string
}

// NewClient creates an initialized ImageKit client.
func NewClient(cfg Config) (*Client, error) {
	if strings.TrimSpace(cfg.PublicKey) == "" {
		return nil, errors.New("ImageKit public key is required")
	}
	if strings.TrimSpace(cfg.PrivateKey) == "" {
		return nil, errors.New("ImageKit private key is required")
	}

	imagekit := imagekitsdk.NewClient(option.WithPrivateKey(cfg.PrivateKey))
	return &Client{
		sdk:         &imagekit,
		publicKey:   cfg.PublicKey,
		privateKey:  cfg.PrivateKey,
		urlEndpoint: strings.TrimRight(cfg.URLEndpoint, "/"),
	}, nil
}

func (c *Client) GenerateUploadAuth(uid, originalFilename string) (UploadAuth, error) {
	folder, err := profileImageFolder(uid)
	if err != nil {
		return UploadAuth{}, err
	}
	if originalFilename == "" || len(originalFilename) > maxOriginalNameBytes ||
		strings.TrimSpace(originalFilename) != originalFilename || !utf8.ValidString(originalFilename) ||
		strings.ContainsAny(originalFilename, `/\`) {
		return UploadAuth{}, ErrInvalidFilename
	}
	extension := strings.ToLower(path.Ext(originalFilename))
	stem := strings.TrimSuffix(originalFilename, path.Ext(originalFilename))
	if stem == "" || stem == "." {
		return UploadAuth{}, ErrInvalidFilename
	}
	switch extension {
	case ".jpg":
		extension = ".jpg"
	case ".jpeg":
		extension = ".jpg"
	case ".png":
		extension = ".png"
	case ".webp":
		extension = ".webp"
	default:
		return UploadAuth{}, ErrUnsupportedExtension
	}

	fileName := "avatar-" + uuid.NewString() + extension
	issuedAt := time.Now().UTC().Unix()
	expiresAt := issuedAt + int64(UploadAuthTTL/time.Second)
	claims := jwt.MapClaims{
		"fileName":          fileName,
		"folder":            folder,
		"useUniqueFileName": "false",
		"overwriteFile":     "false",
		"isPrivateFile":     "false",
		"isPublished":       "true",
		"checks":            UploadChecks,
		"iat":               issuedAt,
		"exp":               expiresAt,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = c.publicKey
	signed, err := token.SignedString([]byte(c.privateKey))
	if err != nil {
		return UploadAuth{}, fmt.Errorf("sign ImageKit upload token: %w", err)
	}

	return UploadAuth{
		Token:     signed,
		ExpiresAt: expiresAt,
		UploadURL: UploadURL,
		UploadPayload: UploadPayload{
			FileName:          fileName,
			Folder:            folder,
			UseUniqueFileName: "false",
			OverwriteFile:     "false",
			IsPrivateFile:     "false",
			IsPublished:       "true",
			Checks:            UploadChecks,
		},
	}, nil
}

func (c *Client) VerifyProfileImageAsset(ctx context.Context, uid, fileID, filePath string) (VerifiedProfileImageAsset, error) {
	if c == nil || c.sdk == nil {
		return VerifiedProfileImageAsset{}, errors.New("ImageKit client is required")
	}
	if strings.TrimSpace(fileID) == "" || strings.TrimSpace(fileID) != fileID {
		return VerifiedProfileImageAsset{}, ErrAssetVerificationFailed
	}
	folder, err := profileImageFolder(uid)
	if err != nil {
		return VerifiedProfileImageAsset{}, err
	}
	if filePath == "" || strings.TrimSpace(filePath) != filePath || !utf8.ValidString(filePath) {
		return VerifiedProfileImageAsset{}, ErrInvalidProfileImagePath
	}
	if path.Clean(filePath) != filePath || path.Dir(filePath) != folder {
		return VerifiedProfileImageAsset{}, ErrInvalidProfileImagePath
	}

	fileName := path.Base(filePath)
	if !strings.HasPrefix(fileName, "avatar-") {
		return VerifiedProfileImageAsset{}, ErrInvalidProfileImagePath
	}
	extension := strings.ToLower(path.Ext(fileName))
	if extension != ".jpg" && extension != ".png" && extension != ".webp" {
		return VerifiedProfileImageAsset{}, ErrInvalidProfileImagePath
	}
	stem := strings.TrimSuffix(fileName, extension)
	avatarID := strings.TrimPrefix(stem, "avatar-")
	parsedAvatarID, err := uuid.Parse(avatarID)
	if err != nil || parsedAvatarID.String() != avatarID {
		return VerifiedProfileImageAsset{}, ErrInvalidProfileImagePath
	}

	asset, err := c.sdk.Files.Get(ctx, fileID)
	if err != nil {
		return VerifiedProfileImageAsset{}, fmt.Errorf("%w: %v", ErrAssetVerificationFailed, err)
	}
	if asset == nil || asset.FileID != fileID || asset.FilePath != filePath {
		return VerifiedProfileImageAsset{}, ErrAssetVerificationFailed
	}
	return VerifiedProfileImageAsset{
		FileID:   asset.FileID,
		FilePath: asset.FilePath,
	}, nil
}

func (c *Client) ProfileImageURL(filePath string) (string, error) {
	if c == nil || c.urlEndpoint == "" {
		return "", errors.New("ImageKit URL endpoint is required")
	}
	if filePath == "" || !strings.HasPrefix(filePath, "/") {
		return "", ErrInvalidProfileImagePath
	}
	return c.urlEndpoint + filePath, nil
}

func profileImageFolder(uid string) (string, error) {
	if uid == "" || strings.TrimSpace(uid) == "" || !utf8.ValidString(uid) {
		return "", ErrInvalidUID
	}
	return profileImageRoot + "/" + base64.RawURLEncoding.EncodeToString([]byte(uid)), nil
}
