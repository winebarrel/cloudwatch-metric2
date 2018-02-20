require 'formula'

class CloudwatchMetric2 < Formula
  VERSION = '0.1.7'

  homepage 'https://github.com/winebarrel/cloudwatch-metric2'
  url "https://github.com/winebarrel/cloudwatch-metric2/releases/download/v#{VERSION}/cloudwatch-metric2-v#{VERSION}-darwin-amd64.gz"
  sha256 '1e19ca4d51b48dcc208eddb574e489eb77a6438c7101d107b7c4d82e16ce7996'
  version VERSION
  head 'https://github.com/winebarrel/cloudwatch-metric2.git', :branch => 'master'

  def install
    system "mv cloudwatch-metric2-v#{VERSION}-darwin-amd64 cloudwatch-metric2"
    bin.install 'cloudwatch-metric2'
  end
end
