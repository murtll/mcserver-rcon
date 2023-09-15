package util

import "net/url"

func SetQueryParam(u *url.URL, key, value string) {
	q := u.Query()
	q.Set(key, value)
	u.RawQuery = q.Encode()
}
