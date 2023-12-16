class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"
  url "https://github.com/deggja/netfetch/releases/download/3.0.0/netfetch_3.0.0_darwin_amd64.tar.gz"
  sha256 "241cc036ed026c04127969b47d82e6cc74f4a62413ca24376424dfeec46f1df3"

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "--version"
  end
end