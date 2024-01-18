class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"
  url "https://github.com/deggja/netfetch/releases/download/3.2.4/netfetch_3.2.4_darwin_amd64.tar.gz"
  sha256 "2e71f0e3fc1fb5c5c8752a500c8f6dfdec432ff26a45e74c4ef939b2eac72d14"

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "--version"
  end
end