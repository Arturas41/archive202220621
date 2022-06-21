package sources

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lxc/distrobuilder/shared"
)

type rockylinux struct {
	commonRHEL

	fname        string
	majorVersion string
}

// Run downloads the tarball and unpacks it.
func (s *rockylinux) Run() error {
	var err error

	s.majorVersion = strings.Split(s.definition.Image.Release, ".")[0]
	baseURL := fmt.Sprintf("%s/%s/isos/%s/", s.definition.Source.URL,
		strings.ToLower(s.definition.Image.Release),
		s.definition.Image.ArchitectureMapped)

	s.fname, err = s.getRelease(s.definition.Source.URL, s.definition.Image.Release,
		s.definition.Source.Variant, s.definition.Image.ArchitectureMapped)
	if err != nil {
		return fmt.Errorf("Failed to get release: %w", err)
	}

	fpath := s.getTargetDir()

	url, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("Failed to parse URL %q: %w", baseURL, err)
	}

	checksumFile := ""
	if !s.definition.Source.SkipVerification {
		// Force gpg checks when using http
		if url.Scheme != "https" {
			if len(s.definition.Source.Keys) == 0 {
				return errors.New("GPG keys are required if downloading from HTTP")
			}

			checksumFile = "CHECKSUM"

			_, err := s.DownloadHash(s.definition.Image, baseURL+checksumFile, "", nil)
			if err != nil {
				return fmt.Errorf("Failed to download %q: %w", baseURL+checksumFile, err)
			}
		}
	}

	_, err = s.DownloadHash(s.definition.Image, baseURL+s.fname, checksumFile, sha256.New())
	if err != nil {
		return fmt.Errorf("Failed to download %q: %w", baseURL+s.fname, err)
	}

	s.logger.WithField("file", filepath.Join(fpath, s.fname)).Info("Unpacking ISO")

	err = s.unpackISO(filepath.Join(fpath, s.fname), s.rootfsDir, s.isoRunner)
	if err != nil {
		return fmt.Errorf("Failed to unpack ISO: %w", err)
	}

	return nil
}

func (s *rockylinux) isoRunner(gpgKeysPath string) error {
	err := shared.RunScript(s.ctx, fmt.Sprintf(`#!/bin/sh
set -eux
GPG_KEYS="%s"
RELEASE="%s"
# Create required files
touch /etc/mtab /etc/fstab
yum_args=""
mkdir -p /etc/yum.repos.d
if [ -d /mnt/cdrom ]; then
	# Install initial package set
	cd /mnt/cdrom/Packages
	rpm -ivh --nodeps $(ls yum-*.rpm | head -n1)
	# Add cdrom repo
	cat <<- EOF > /etc/yum.repos.d/cdrom.repo
[cdrom]
name=Install CD-ROM
baseurl=file:///mnt/cdrom
enabled=0
EOF
	if [ -n "${GPG_KEYS}" ]; then
		echo gpgcheck=1 >> /etc/yum.repos.d/cdrom.repo
		echo gpgkey=${GPG_KEYS} >> /etc/yum.repos.d/cdrom.repo
	else
		echo gpgcheck=0 >> /etc/yum.repos.d/cdrom.repo
	fi
	yum_args="--disablerepo=* --enablerepo=cdrom"
	yum ${yum_args} -y --releasever="${RELEASE}" reinstall yum
else
	if ! [ -f /etc/pki/rpm-gpg/RPM-GPG-KEY-rockyofficial ]; then
		mkdir -p /etc/pki/rpm-gpg
		cat <<- "EOF" > /etc/pki/rpm-gpg/RPM-GPG-KEY-rockyofficial
-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v2.0.22 (GNU/Linux)
mQINBGAofzYBEAC6yS1azw6f3wmaVd//3aSy6O2c9+jeetulRQvg2LvhRRS1eNqp
/x9tbBhfohu/tlDkGpYHV7diePgMml9SZDy1sKlI3tDhx6GZ3xwF0fd1vWBZpmNk
D9gRkUmYBeLotmcXQZ8ZpWLicosFtDpJEYpLUhuIgTKwt4gxJrHvkWsGQiBkJxKD
u3/RlL4IYA3Ot9iuCBflc91EyAw1Yj0gKcDzbOqjvlGtS3ASXgxPqSfU0uLC9USF
uKDnP2tcnlKKGfj0u6VkqISliSuRAzjlKho9Meond+mMIFOTT6qp4xyu+9Dj3IjZ
IC6rBXRU3xi8z0qYptoFZ6hx70NV5u+0XUzDMXdjQ5S859RYJKijiwmfMC7gZQAf
OkdOcicNzen/TwD/slhiCDssHBNEe86Wwu5kmDoCri7GJlYOlWU42Xi0o1JkVltN
D8ZId+EBDIms7ugSwGOVSxyZs43q2IAfFYCRtyKHFlgHBRe9/KTWPUrnsfKxGJgC
Do3Yb63/IYTvfTJptVfhQtL1AhEAeF1I+buVoJRmBEyYKD9BdU4xQN39VrZKziO3
hDIGng/eK6PaPhUdq6XqvmnsZ2h+KVbyoj4cTo2gKCB2XA7O2HLQsuGduHzYKNjf
QR9j0djjwTrsvGvzfEzchP19723vYf7GdcLvqtPqzpxSX2FNARpCGXBw9wARAQAB
tDNSZWxlYXNlIEVuZ2luZWVyaW5nIDxpbmZyYXN0cnVjdHVyZUByb2NreWxpbnV4
Lm9yZz6JAk4EEwEIADgWIQRwUcRwqSn0VM6+N7cVr12sbXRaYAUCYCh/NgIbDwUL
CQgHAgYVCgkICwIEFgIDAQIeAQIXgAAKCRAVr12sbXRaYLFmEACSMvoO1FDdyAbu
1m6xEzDhs7FgnZeQNzLZECv2j+ggFSJXezlNVOZ5I1I8umBan2ywfKQD8M+IjmrW
k9/7h9i54t8RS/RN7KNo7ECGnKXqXDPzBBTs1Gwo1WzltAoaDKUfXqQ4oJ4aCP/q
/XPVWEzgpJO1XEezvCq8VXisutyDiXEjjMIeBczxb1hbamQX+jLTIQ1MDJ4Zo1YP
zlUqrHW434XC2b1/WbSaylq8Wk9cksca5J+g3FqTlgiWozyy0uxygIRjb6iTzKXk
V7SYxeXp3hNTuoUgiFkjh5/0yKWCwx7aQqlHar9GjpxmBDAO0kzOlgtTw//EqTwR
KnYZLig9FW0PhwvZJUigr0cvs/XXTTb77z/i/dfHkrjVTTYenNyXogPtTtSyxqca
61fbPf0B/S3N43PW8URXBRS0sykpX4SxKu+PwKCqf+OJ7hMEVAapqzTt1q9T7zyB
QwvCVx8s7WWvXbs2d6ZUrArklgjHoHQcdxJKdhuRmD34AuXWCLW+gH8rJWZpuNl3
+WsPZX4PvjKDgMw6YMcV7zhWX6c0SevKtzt7WP3XoKDuPhK1PMGJQqQ7spegGB+5
DZvsJS48Ip0S45Qfmj82ibXaCBJHTNZE8Zs+rdTjQ9DS5qvzRA1sRA1dBb/7OLYE
JmeWf4VZyebm+gc50szsg6Ut2yT8hw==
=AiP8
-----END PGP PUBLIC KEY BLOCK-----
EOF
	fi
	cat <<- "EOF" > /etc/yum.repos.d/Rocky-BaseOS.repo
[BaseOS]
name=Rocky-$releasever - Base
mirrorlist=http://mirrors.rockylinux.org/mirrorlist?arch=$basearch&repo=BaseOS-8
gpgcheck=1
enabled=1
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-rockyofficial
EOF
	# Use dnf in the boot iso since yum isn't available
	alias yum=dnf
fi
pkgs="basesystem Rocky-release yum"
# Create a minimal rootfs
mkdir /rootfs
yum ${yum_args} --installroot=/rootfs -y --releasever="${RELEASE}" --skip-broken install ${pkgs}
rm -rf /rootfs/var/cache/yum
`, gpgKeysPath, s.majorVersion))
	if err != nil {
		return fmt.Errorf("Failed to run ISO script: %w", err)
	}

	return nil
}

func (s *rockylinux) getRelease(URL, release, variant, arch string) (string, error) {
	u := URL + path.Join("/", strings.ToLower(release), "isos", arch)

	resp, err := http.Get(u)
	if err != nil {
		return "", fmt.Errorf("Failed to GET %q: %w", u, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read body: %w", err)
	}

	re := s.getRegexes(arch, variant, release)

	for _, r := range re {
		matches := r.FindAllString(string(body), -1)
		if len(matches) > 0 {
			return matches[len(matches)-1], nil
		}
	}

	return "", errors.New("Failed to find release")
}

func (s *rockylinux) getRegexes(arch string, variant string, release string) []*regexp.Regexp {
	releaseFields := strings.Split(release, ".")

	var re []string
	switch len(releaseFields) {
	case 1:
		re = append(re, fmt.Sprintf("Rocky-%s(.\\d+)*-%s-(?i:%s).iso",
			releaseFields[0], arch, variant))
	case 2:
		re = append(re, fmt.Sprintf("Rocky-%s.%s-%s-(?i:%s).iso",
			releaseFields[0], releaseFields[1], arch, variant))
	}

	regexes := make([]*regexp.Regexp, len(re))

	for i, r := range re {
		regexes[i] = regexp.MustCompile(r)
	}

	return regexes
}
