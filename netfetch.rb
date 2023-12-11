class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/0.0.65/netfetch_0.0.65_darwin_amd64.tar.gz"
    sha256 "ea8ca9b57223d69eb2db593080fa4cc1f9abfe1908e2aed3253fae4fd8614fcd"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end