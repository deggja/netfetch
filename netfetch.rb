class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"
  url "https://github.com/deggja/netfetch/releases/download/2.1.0/netfetch_2.1.0_darwin_amd64.tar.gz"
  sha256 "7b0d01392469ee9c1c035921f0420692220689f0e6654ee11bd22d4e7c05245a"

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "--version"
  end
end