class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"
  url "https://github.com/deggja/netfetch/releases/download/3.0.8/netfetch_3.0.8_darwin_amd64.tar.gz"
  sha256 "fad340cae0ed1d06c65df7a0864594b3c5e272bf9b1d59cbffda3f027197fef9"

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "--version"
  end
end