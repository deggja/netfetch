class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"
  url "https://github.com/deggja/netfetch/releases/download/3.0.4/netfetch_3.0.4_darwin_amd64.tar.gz"
  sha256 "18a405006c6d1adabb74a821eeecb23deb5952d387d8708f4acc923736a5503a"

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "--version"
  end
end