package middleware

import (
	"strings"
	"testing"
)

func TestRedactPasswords(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{
			name: "登录：密码脱敏，邮箱保留",
			in:   `{"email":"a@b.com","password":"P@ssw0rd-Example"}`,
			want: `{"email":"a@b.com","password":"***"}`,
		},
		{
			name: "注册：password 脱敏，其余字段保留",
			in:   `{"username":"x","password":"test-pw-Ab3","email":"a@b.com","code":"123456"}`,
			want: `{"username":"x","password":"***","email":"a@b.com","code":"123456"}`,
		},
		{
			name: "改密：old_password / new_password 都脱敏",
			in:   `{"old_password":"test-pw-Ab3","new_password":"new-pw-Ab3"}`,
			want: `{"old_password":"***","new_password":"***"}`,
		},
		{
			name: "不含密码字段的请求体原样保留",
			in:   `{"scene":"avatar","sort":1}`,
			want: `{"scene":"avatar","sort":1}`,
		},
	}
	for _, c := range cases {
		if got := redactPasswords(c.in); got != c.want {
			t.Errorf("%s: redactPasswords(%q) = %q, want %q", c.name, c.in, got, c.want)
		}
	}
}

// TestRedactPasswordsNoLeak 确保脱敏后日志里搜不到明文密码
func TestRedactPasswordsNoLeak(t *testing.T) {
	for _, pw := range []string{"P@ssw0rd-Example", "test-pw-Ab3", "new-pw-Ab3"} {
		body := `{"email":"a@b.com","password":"` + pw + `","new_password":"` + pw + `"}`
		got := redactPasswords(body)
		if strings.Contains(got, pw) {
			t.Fatalf("脱敏后仍残留明文密码 %q: %s", pw, got)
		}
	}
}
