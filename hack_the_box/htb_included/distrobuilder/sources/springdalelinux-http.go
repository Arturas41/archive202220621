package sources

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lxc/distrobuilder/shared"
)

type springdalelinux struct {
	commonRHEL

	fname        string
	majorVersion string
}

// Run downloads the tarball and unpacks it.
func (s *springdalelinux) Run() error {
	s.majorVersion = strings.Split(s.definition.Image.Release, ".")[0]

	// Example: http://puias.princeton.edu/data/puias/8.3/x86_64/os/images/boot.iso
	baseURL := fmt.Sprintf("%s/%s/%s/os/images/", s.definition.Source.URL,
		strings.ToLower(s.definition.Image.Release),
		s.definition.Image.ArchitectureMapped)
	s.fname = "boot.iso"

	fpath := s.getTargetDir()

	_, err := s.DownloadHash(s.definition.Image, baseURL+s.fname, "", nil)
	if err != nil {
		return fmt.Errorf("Error downloading %q: %w", baseURL+s.fname, err)
	}

	s.logger.WithField("file", filepath.Join(fpath, s.fname)).Info("Unpacking ISO")

	err = s.unpackISO(filepath.Join(fpath, s.fname), s.rootfsDir, s.isoRunner)
	if err != nil {
		return fmt.Errorf("Failed to unpack %q: %w", filepath.Join(fpath, s.fname), err)
	}

	return nil
}

func (s *springdalelinux) isoRunner(gpgKeysPath string) error {
	err := shared.RunScript(s.ctx, fmt.Sprintf(`#!/bin/sh
set -eux

GPG_KEYS="%s"

# Create required files
touch /etc/mtab /etc/fstab

yum_args=""
mkdir -p /etc/yum.repos.d

if [ -d /mnt/cdrom ]; then
	# Install initial package set
	cd /mnt/cdrom/Packages
	rpm -ivh --nodeps $(ls rpm-*.rpm | head -n1)
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
	yum ${yum_args} -y reinstall yum
else
	if ! [ -f /etc/pki/rpm-gpg/RPM-GPG-KEY-springdale ]; then
		mkdir -p /etc/pki/rpm-gpg
		cat <<- "EOF" > /etc/pki/rpm-gpg/RPM-GPG-KEY-springdale
-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v2.0.14 (GNU/Linux)

mQILBEzhkmYBEAD1+atknNh1ufb0QZNKl6uMnOnAL2bk8v6VmUjVTNTJ6JBXdREX
omVRlEDv1Jw9/11tDEI0JzZieqogahJ3GDqENXhLlfpgev+YHES6jrH8XEYkin/T
wxuUyKvqoAk60wqC6Uv/TkHpVOYP9WHtNbybl3Pu8f2tzqdrE20HmQgEzVjTyHEw
33WeAXWlhb/rIMcJ/edk0/hQ0Yu1NR34g/R/1azfz4uyx0qtrTDvaS7qxX5ZSUYZ
tHzS8JjFpzRHI3A++Jv6yXMdzekn54TYob+DeKoRsdtfI0KryOv+92tcw4yM35CK
d4apq8t35v56wsRXhwVZWEYn3BfIsdDoCF4XXH8nx0B1KJii6uYo1iiyIU10cT+g
YCPVIzT+B/FqAmltx3VJg87sZ6QVFp3oFtIOlu14NQBv3wqUneCNK88EZAxi/kqy
kFLq4NfWpn2dVZYRlJddWTXZ+qqaYQ6aI9HyGx5/QTI1hsxcOykrm8woZvptYeRQ
IfkDFPgqTmdTqFaL5TQc43FvOLss9kBJ2FSbi3WpRIJW0/AQGc9/97SUnC0T075O
R3cK+dQvI3QV3/p1UEowXlSHDFc5CrqKT11zMG9KE4UGm08yDbz0d5IbH0Cpzuku
vidjBlsaQ/+OldjSFjHhCFLz55lLya6eMpPaScEKqoWNYzuNv9kDgVWXNwAGKbRN
UFVJQVMgTGludXggRGlzdHJpYnV0aW9uIChSUE0gSW50ZWdyaXR5IFNpZ25hdHVy
ZSkgPHB1aWFzQG1hdGgucHJpbmNldG9uLmVkdT6JAjQEEwECAB4FAkzhkmYCGwMG
CwkIBwMCAxUCAwMWAgECHgECF4AACgkQFs/DM0GkCUi22xAAtoeFPRpYAoaq6+Ru
nRX5GCDQl6DlOtVxLVclNZzGpnw8Extid+AOLqDXcfncyf04YhlEHj4misz/rDCI
a5bRWoNjPHAgHpzCX7+I6pNr1hY9SOW94BEdng9IGGK0XhFBzflmySLZEC9E2ZYe
RgWKJcDbyM9sDc2g440ICkn8DOWTvKMcQ7f0AzYtARXfmAEMqgqzNV+0wDJmdEHY
7rif51U8bCOKns/UFKSA3WqUKhn5v2xo4OVqkm+bVG3z04KRAIUWIZIK8RHEp6wk
clls8afYSufJmmUeczbE/wDqEgMSE3qGlcQRxTO3EMORb3nwWo5QAA4I/QPFrFoC
QZQbaLNOx8P7dnfDoarJrUPYBBUFmMKvHUnSwv696QZhz70RvgjTHcbSyrnmE76C
/XU0zUpeWN6FEb77zA1pIlqVf4hqRs+PCaG2sytBQVYEpYgGnoUPSWIT4a6NJtn+
WwJHOqRYrGGTM0Z6V7IgMAkqiwEECn5eDUXYhqwUsyuVkbeOBWTc6nhsPIH3QC86
sL0X/1hztP5sDoCne6SY3X3IyglvApvsKn0TOcVCvbYNhbg2bfPdvfmtSAV3/iMU
yPw2JgfcvLeF1tiMQ7i5PfgyOn0Y3/lZjcclYHa1P5PEoTCA2lU7jLm0lmgIrh4N
vYD0DGRZTGNEJYUDZFgRyynIOxk=
=mKoc
-----END PGP PUBLIC KEY BLOCK-----
EOF
	fi

	. /etc/os-release

	if [[ ${VERSION_ID:0:1} == "7" ]]
	then
		fname="$(curl -s http://springdale.princeton.edu/data/springdale/7/x86_64/os/Packages/ | grep -Eo 'yum-[[:digit:]].+\.noarch\.rpm')"

		if [ -n "${fname}" ]; then
			wget http://springdale.princeton.edu/data/springdale/7/x86_64/os/Packages/"${fname}"
			rpm -ivh --nodeps yum*.rpm
		fi
	fi


	if [[ ${VERSION_ID:0:1} == "8" ]]
	then
	cat <<- "EOF" > /etc/yum.repos.d/Springdale-Base.repo
[sdl8-baseos]
name=Springdale core Base $releasever - $basearch
mirrorlist=http://springdale.princeton.edu/data/springdale/$releasever/$basearch/os/BaseOS/mirrorlist
#baseurl=http://springdale.princeton.edu/data/springdale/$releasever/$basearch/os/BaseOS
gpgcheck=1
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-springdale

[sdl8-baseos-updates]
name=Springdale core Updates $releasever - $basearch
mirrorlist=http://springdale.princeton.edu/data/springdale/updates/$releasever/BaseOS/$basearch/mirrorlist
#baseurl=http://springdale.princeton.edu/data/springdale/updates/$releasever/BaseOS/$basearch
gpgcheck=1
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-springdale
EOF
	# Use dnf in the boot iso since yum isn't available
	alias yum=dnf
	fi

	if [[ ${VERSION_ID:0:1} == "7" ]]
	then
	cat <<- "EOF" > /etc/yum.repos.d/Springdale-Base.repo
[core]
name=Springdale core Base $releasever - $basearch
#baseurl=file:///springdale/$releasever/$basearch/os
#mirrorlist=http://mirror.math.princeton.edu/pub/springdale/puias/$releasever/$basearch/os/mirrorlist
baseurl=http://springdale.princeton.edu/data/springdale/$releasever/$basearch/os
gpgcheck=1
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-springdale

[updates]
name=Springdale core Updates $releasever - $basearch
#baseurl=file:///springdale/updates/$releasever/en/os/$basearch
#mirrorlist=http://mirror.math.princeton.edu/pub/springdale/puias/updates/$releasever/en/os/$basearch/mirrorlist
baseurl=http://springdale.princeton.edu/data/springdale/updates/$releasever/en/os/$basearch
gpgcheck=1
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-springdale
EOF
	fi
fi

pkgs="basesystem springdale-release yum"

# Create a minimal rootfs
mkdir /rootfs
echo "install rootfs"
yum ${yum_args} --installroot=/rootfs -y --releasever=%s --skip-broken install ${pkgs}
rm -rf /rootfs/var/cache/yum
`, gpgKeysPath, s.majorVersion))
	if err != nil {
		return fmt.Errorf("Failed to run ISO script: %w", err)
	}

	return nil
}
