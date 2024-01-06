class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"
  url "https://github.com/deggja/netfetch/releases/download/3.0.6/netfetch_3.0.6_darwin_amd64.tar.gz"
  sha256 "c6557dcbb80fa53d871402e7b933d098bd5a8ae6ef0780ba252ebbfbe763a4d6"

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "--version"
  end
end