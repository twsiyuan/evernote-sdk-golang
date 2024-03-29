// Code generated by Thrift Compiler (0.16.0). DO NOT EDIT.

package edam

import (
	"bytes"
	"context"
	"fmt"
	"time"
	thrift "github.com/apache/thrift/lib/go/thrift"

)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = context.Background
var _ = time.Now
var _ = bytes.Equal

const CLASSIFICATION_RECIPE_USER_NON_RECIPE = "000"
const CLASSIFICATION_RECIPE_USER_RECIPE = "001"
const CLASSIFICATION_RECIPE_SERVICE_RECIPE = "002"
const EDAM_NOTE_SOURCE_WEB_CLIP = "web.clip"
const EDAM_NOTE_SOURCE_WEB_CLIP_SIMPLIFIED = "Clearly"
const EDAM_NOTE_SOURCE_MAIL_CLIP = "mail.clip"
const EDAM_NOTE_SOURCE_MAIL_SMTP_GATEWAY = "mail.smtp"

func init() {
}

