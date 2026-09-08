package service

import (
	"testing"
	"time"
)

// extFromMime 决定上传文件扩展名白名单，是上传安全的关键一环（不信任前端 Content-Type）
func TestExtFromMime(t *testing.T) {
	cases := []struct {
		mime string
		want string
		ok   bool
	}{
		{"image/jpeg", ".jpg", true},
		{"image/png", ".png", true},
		{"image/gif", ".gif", true},
		{"image/webp", ".webp", true},
		{"image/svg+xml", "", false},
		{"application/x-msdownload", "", false},
		{"text/html", "", false},
		{"", "", false},
	}
	for _, c := range cases {
		got, ok := extFromMime(c.mime)
		if got != c.want || ok != c.ok {
			t.Errorf("extFromMime(%q) = (%q,%v)，期望 (%q,%v)", c.mime, got, ok, c.want, c.ok)
		}
	}
}

func TestRandomHex(t *testing.T) {
	a, b := randomHex(16), randomHex(16)
	if len(a) != 32 || len(b) != 32 {
		t.Fatalf("randomHex(16) 应返回 32 个 hex 字符，实际 %d/%d", len(a), len(b))
	}
	if a == b {
		t.Error("randomHex 两次结果相同，随机性异常")
	}
}

func TestSecondsUntilMonthEnd(t *testing.T) {
	// 月末时刻：剩余时间应为正且不超过 1 天（1 月 31 日 23:59:59 → 次月 1 日）
	d := secondsUntilMonthEnd(time.Date(2026, 1, 31, 23, 59, 59, 0, time.Local))
	if d <= 0 || d > 24*time.Hour {
		t.Errorf("月末剩余秒数异常: %v", d)
	}
}

func TestMd5Hash(t *testing.T) {
	if md5Hash("a") == md5Hash("b") {
		t.Error("md5Hash 不同输入得到相同结果")
	}
	if h := md5Hash("abc"); len(h) != 32 {
		t.Errorf("md5Hash 长度异常: %d", len(h))
	}
}
