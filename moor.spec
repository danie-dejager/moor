Name:    moor
Summary: Simple UTF-8 pager with sensible defaults
Version: 2.10.5
Release: 1%{?dist}
License: BSD-2-Clause
URL:     https://github.com/walles/moor
Source0: https://github.com/walles/moor/archive/refs/tags/v%{version}.tar.gz

%define debug_package %{nil}

BuildRequires: curl
BuildRequires: gcc
BuildRequires: make
BuildRequires: gzip
BuildRequires: golang
BuildRequires: git

%description
Moor is a pager for UTF-8 encoded text. It reads and displays
text from files or from pipelines. It is designed to work out of
the box with sensible defaults, without requiring user configuration.

%prep
%setup -q

%build
GO111MODULE=on go build -v -trimpath -modcacherw \
   -ldflags="-s -w -X main.versionString=%{version}" \
   -o %{name} ./cmd/%{name}

strip --strip-all %{name}
gzip %{name}.1

%install
mkdir -p %{buildroot}%{_bindir}
mkdir -p %{buildroot}%{_mandir}/man1
install -m 755 %{name} %{buildroot}%{_bindir}
install -m 644 %{name}.1.gz %{buildroot}%{_mandir}/man1

%files
%doc README.md
%license LICENSE
%{_bindir}/%{name}
%{_mandir}/man1/%{name}.1.gz

%changelog
* Sun Feb 22 2026 - Danie de Jager <danie.dejager@gmail.com> - 2.10.5-1
- Handle paging non-seekable files
* Mon Feb 9 2026 - Danie de Jager <danie.dejager@gmail.com> - 2.10.4-1
- Fix bug pressing "n" at the bottom of the input.
* Wed Jan 28 2026 - Danie de Jager <danie.dejager@gmail.com> - 2.10.3-1
- Fix two crashes
* Mon Jan 19 2026 - Danie de Jager <danie.dejager@gmail.com> - 2.10.2-1
* Fri Jan 2 2026 - Danie de Jager <danie.dejager@gmail.com> - 2.10.1-1
* Tue Dec 17 2025 - Danie de Jager <danie.dejager@gmail.com> - 2.9.6-1
- Various improvements
* Tue Dec 9 2025 - Danie de Jager <danie.dejager@gmail.com> - 2.9.5-1
- Fix non-working case insensitive search
* Sat Dec 6 2025 - Danie de Jager <danie.dejager@gmail.com> - 2.9.4-1
* Sun Nov 30 2025 - Danie de Jager <danie.dejager@gmail.com> - 2.9.3-1
- Search performance improvements
* Sat Nov 15 2025 - Danie de Jager <danie.dejager@gmail.com> - 2.9.2-1
- 4x Speed improvement
* Sat Nov 15 2025 - Danie de Jager <danie.dejager@gmail.com> - 2.9.1-1
* Fri Nov 14 2025 - Danie de Jager <danie.dejager@gmail.com> - 2.9.0-1
- Add persistent search history
* Sun Nov 9 2025 - Danie de Jager <danie.dejager@gmail.com> - 2.8.2-1
- Fix search hit line highlighting with word wrapping enabled.
- Show keyboard help while searching
- Inform user when toggling word wrapping
* Mon Nov 3 2025 - Danie de Jager <danie.dejager@gmail.com> - 2.8.1-1
- searching could sometimes scroll right
- before this release, in Kitty and some other terminals, mouse selection didn't work while content was still loading. This should now be fixed.
- Add LESSSECURE=1 secure mode for systemctl
* Sun Oct 26 2025 - Danie de Jager <danie.dejager@gmail.com> - 2.7.1-1
- Improve --terminal-fg with terminal bg images
- Use PAGER_LABEL env var to label stdin
* Sun Oct 19 2025 - Danie de Jager <danie.dejager@gmail.com> - 2.6.1-1
- Support `#` in URLs 
- Add QuitIfOneScreen and NoLineNumbers to embed API
* Wed Oct 15 2025 - Danie de Jager <danie.dejager@gmail.com> - 2.5.2-1
- Center search hits vertically
* Sun Oct 12 2025 - Danie de Jager <danie.dejager@gmail.com> - 2.5.1-1
- Backwards searching now scrolls sideways as needed
* Tue Oct 7 2025 - Danie de Jager <danie.dejager@gmail.com> - 2.5.0-1
- Make forward search find sideways matches
* Fri Oct 3 2025 - Danie de Jager <danie.dejager@gmail.com>- 2.4.1-1
- Match less' behavior with piped stdin
* Mon Sep 29 2025 - Danie de Jager <danie.dejager@gmail.com>- 2.4.0-1
- Default tab size to 8 to be like less
* Thu Sep 25 2025 - Danie de Jager <danie.dejager@gmail.com>- 2.3.0-1
- Support opening multiple files
* Wed Sep 24 2025 - Danie de Jager <danie.dejager@gmail.com>- 2.2.1-1
- Provide line highlighting in more cases
* Sun Sep 21 2025 - Danie de Jager <danie.dejager@gmail.com>- 2.2.0-1
- Highlight lines with search hits
* Sat Sep 13 2025 - Danie de Jager <danie.dejager@gmail.com>- 2.1.1-1
- Working scroll + select in Windows Terminal
* Sun Aug 31 2025 - Danie de Jager <danie.dejager@gmail.com>- 2.1.0-1
- Accept `-` to mean "read from `stdin`"
* Tue Aug 26 2025 - Danie de Jager <danie.dejager@gmail.com>- 2.0.5-1
- Fixed a crash related to intermittent problem related to scrolling around the switch from line numbers 999 to 1000.
- Mac keyboards can now press option-arrow to scroll sideways one column at a time.
* Sat Aug 23 2025 - Danie de Jager <danie.dejager@gmail.com>- 2.0.4-2
- Update license
