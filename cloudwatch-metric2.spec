%define  debug_package %{nil}

Name:     cloudwatch-metric2
Version:  0.1.6
Release:  1%{?dist}
Summary:  a tool to get AWS CloudWatch metrics.

Group:    Development/Tools
License:  MIT
URL:      https://github.com/winebarrel/cloudwatch-metric2
Source0:  %{name}.tar.gz
# https://github.com/winebarrel/cloudwatch-metric2/releases/download/v%{version}/cloudwatch-metric2_%{version}.tar.gz

%description
a tool to get AWS CloudWatch metrics.

%prep
%setup -q -n src

%build
make

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}/usr/bin
install -m 755 cloudwatch-metric2 %{buildroot}/usr/bin/

%files
%defattr(755,root,root,-)
/usr/bin/cloudwatch-metric2
