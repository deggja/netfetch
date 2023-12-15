class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"
  url "https://github.com/deggja/netfetch/releases/download/2.0.0/netfetch_2.0.0_darwin_amd64.tar.gz"
  sha256 "0aba03542eb6137382be8fe61d5801734e2771250ae8f0355add438b26aae2c3"

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "--version"
  end
end