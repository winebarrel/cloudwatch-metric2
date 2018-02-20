require 'formula'

class CloudwatchMetric2 < Formula
  VERSION = '0.1.7'

  homepage 'https://github.com/winebarrel/cloudwatch-metric2'
  url "https://github.com/winebarrel/cloudwatch-metric2/releases/download/v#{VERSION}/cloudwatch-metric2-v#{VERSION}-darwin-amd64.gz"
  sha256 '709bed8302310afd2f8d3d60342216265188b2fa42ffece7dbe3eb958a81458a'
  version VERSION
  head 'https://github.com/winebarrel/cloudwatch-metric2.git', :branch => 'master'

  def install
    system "mv cloudwatch-metric2-v#{VERSION}-darwin-amd64 cloudwatch-metric2"
    bin.install 'cloudwatch-metric2'
  end
end
