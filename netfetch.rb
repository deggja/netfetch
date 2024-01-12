class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"
  url "https://github.com/deggja/netfetch/releases/download/3.0.7/netfetch_3.0.7_darwin_amd64.tar.gz"
  sha256 "d7c26d63afb42d525fde196e4b8a3ccfd4a04fd3a1682cb4083ee18b8bad53bc"

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "--version"
  end
end