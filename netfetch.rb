class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"
  url "https://github.com/deggja/netfetch/releases/download/1.5.0/netfetch_1.5.0_darwin_amd64.tar.gz"
  sha256 "8bb003855ade22f2971ce3f988c0855b9e3b703e2214b6bfec783dc69e010057"

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "--version"
  end
end