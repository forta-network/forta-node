Summary: Forta scanner node CLI
Name: forta
Version: SEMVER
Release: 1
License: FIXME
Group: System
Packager: Forta Protocol

BuildRoot: %{buildroot}

Provides: forta
Requires: systemd

%description
Forta scanner node CLI

%install

%files
/usr/lib/systemd/system/forta.service

%post
curl ARTIFACTS_URL/forta-REVISION -o forta -s
install -m 0755 forta %{_bindir}/forta

%postun
rm -f %{_bindir}/forta
rm -f /etc/systemd/user/forta.service
