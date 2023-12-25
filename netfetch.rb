class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"
  url "https://github.com/deggja/netfetch/releases/download/3.0.2/netfetch_3.0.2_darwin_amd64.tar.gz"
  sha256 "35e5821c1d368d1bb91a45d13ef8b0af6f30e3a062c647ecdfebe7e6ea0fa690"

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "--version"
  end
end