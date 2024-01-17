class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"
  url "https://github.com/deggja/netfetch/releases/download/3.2.0/netfetch_3.2.0_darwin_amd64.tar.gz"
  sha256 "aff4377c3cf668725d557ba9bce0b2d7e2ca83cc76b950cd15293fd80de6e2f3"

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "--version"
  end
end