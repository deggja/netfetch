class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"
  url "https://github.com/deggja/netfetch/releases/download/3.2.3/netfetch_3.2.3_darwin_amd64.tar.gz"
  sha256 "2a117d945c2be4b6561863029148f3832a6c82d44d5f478143e5e6043cd736ff"

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "--version"
  end
end