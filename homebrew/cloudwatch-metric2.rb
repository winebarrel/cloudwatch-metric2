require 'formula'

class CloudwatchMetric2 < Formula
  VERSION = '0.1.5'

  homepage 'https://github.com/winebarrel/cloudwatch-metric2'
  url "https://github.com/winebarrel/cloudwatch-metric2/releases/download/v#{VERSION}/cloudwatch-metric2-v#{VERSION}-darwin-amd64.gz"
  sha256 '9b5ced0d44fcf7f9b13dcf5e1f0a12c425a287f495032c372fa51283410d5425'
  version VERSION
  head 'https://github.com/winebarrel/cloudwatch-metric2.git', :branch => 'master'

  def install
    system "mv cloudwatch-metric2-v#{VERSION}-darwin-amd64 cloudwatch-metric2"
    bin.install 'cloudwatch-metric2'
  end
end
