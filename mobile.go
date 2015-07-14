// Package mobile provides methods that determine whether a request by a client is coming
// from a mobile, tablet or normal device. This middleware was inspired by the spring-mobile
// project https://github.com/spring-projects/spring-mobile
package mobile

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

const (
	DefaultKey  = "github.com/floresj/go-contrib-mobile"
	Android     = "android"
	Mobile      = "mobile"
	Ipad        = "ipad"
	Iphone      = "iphone"
	Ipod        = "ipod"
	Wap         = "wap"
	XwapProfile = "X-Wap-Profile"
	Profile     = "Profile"
)

var MobileUserAgentPrefixes = []string{
	"w3c ", "w3c-", "acs-", "alav", "alca", "amoi", "audi", "avan", "benq",
	"bird", "blac", "blaz", "brew", "cell", "cldc", "cmd-", "dang", "doco",
	"eric", "hipt", "htc_", "inno", "ipaq", "ipod", "jigs", "kddi", "keji",
	"leno", "lg-c", "lg-d", "lg-g", "lge-", "lg/u", "maui", "maxo", "midp",
	"mits", "mmef", "mobi", "mot-", "moto", "mwbp", "nec-", "newt", "noki",
	"palm", "pana", "pant", "phil", "play", "port", "prox", "qwap", "sage",
	"sams", "sany", "sch-", "sec-", "send", "seri", "sgh-", "shar", "sie-",
	"siem", "smal", "smar", "sony", "sph-", "symb", "t-mo", "teli", "tim-",
	"tosh", "tsm-", "upg1", "upsi", "vk-v", "voda", "wap-", "wapa", "wapi",
	"wapp", "wapr", "webc", "winw", "winw", "xda ", "xda-"}

var MobileUserAgentKeywords = []string{
	"blackberry", "webos", "ipod", "lge vx", "midp", "maemo", "mmp", "mobile",
	"netfront", "hiptop", "nintendo DS", "novarra", "openweb", "opera mobi",
	"opera mini", "palm", "psp", "phone", "smartphone", "symbian", "up.browser",
	"up.link", "wap", "windows ce"}

var TabletUserAgentKeywords = []string{"ipad", "playbook", "hp-tablet", "kindle"}

type device struct {
	normal   bool
	mobile   bool
	tablet   bool
	platform string
}

type Device interface {
	Normal() bool
	Mobile() bool
	Tablet() bool
	Platform() string
}

func Resolver() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Initializing mobile resolver")
		d := resolveDevice(c.Request.Header)
		log.Println("Setting the Default Key")
		c.Set(DefaultKey, d)
		log.Println("Default key set")
		c.Next()
	}
}

// Reads the Header from a Request and attempts to determine what type of device the user is using.
// Utilizes various checks using the User-Agent,
func resolveDevice(header http.Header) Device {
	agent := strings.ToLower(header.Get("User-Agent"))

	// Check Tablet
	if agent != "" {
		switch {
		case strings.Contains(agent, Android) && !strings.Contains(agent, Mobile):
			return &device{tablet: true}
		case strings.Contains(agent, Ipad):
			return &device{tablet: true}
		case strings.Contains(agent, "silk") && !strings.Contains(agent, Mobile):
			return &device{tablet: true}
		default:
			for _, keyword := range TabletUserAgentKeywords {
				if strings.Contains(agent, keyword) {
					return &device{tablet: true}
				}
			}
		}
	}

	// User Agent Profile detection
	xWapProfile := header.Get(XwapProfile)
	profile := header.Get(Profile)

	if xWapProfile != "" || profile != "" {
		if agent != "" {
			switch {
			case strings.Contains(agent, Android):
				return &device{mobile: true}
			case strings.Contains(agent, Iphone) || strings.Contains(agent, Ipod) || strings.Contains(agent, Ipad):
				return &device{mobile: true}
			default:
				return &device{mobile: true}
			}
		}
	}

	// User Agent Prefix check
	if agent != "" && len(agent) >= 4 {
		prefix := agent[:4]
		for _, uaprefix := range MobileUserAgentPrefixes {
			if strings.Contains(prefix, uaprefix) {
				return &device{mobile: true}
			}
		}
	}

	// Accept Header check
	accept := header.Get("Accept")
	if accept != "" && strings.Contains(accept, Wap) {
		return &device{mobile: true}
	}

	// Check Mobile
	if agent != "" {
		switch {
		case strings.Contains(agent, Android):
			return &device{mobile: true}
		case strings.Contains(agent, Iphone) || strings.Contains(agent, Ipod) || strings.Contains(agent, Ipad):
			return &device{mobile: true}
		default:
			for _, keyword := range MobileUserAgentKeywords {
				if strings.Contains(agent, keyword) {
					return &device{mobile: true}
				}
			}
		}
	}

	// Assume 'normal' if mobile or tablet was not identified
	return &device{normal: true}
}

// Shortcut for retrieving a Device object within a handler
func GetDevice(c *gin.Context) Device {
	return c.MustGet(DefaultKey).(Device)
}

// Returns true if a device is Normal. Normal meaning not a Mobile or Tablet.
func (d *device) Normal() bool {
	return d.normal
}

// Returns true if a device is a mobile device
func (d *device) Mobile() bool {
	return d.mobile
}

// Returns true if a device is a tablet
func (d *device) Tablet() bool {
	return d.tablet
}

// Returns the device platform. Possible values include IOS, Android or Unknown
func (d *device) Platform() string {
	return d.platform
}
