class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/1.1.0/netfetch_1.1.0_darwin_amd64.tar.gz"
    sha256 "0019dfc4b32d63c1392aa264aed2253c1e0c2fb09216f8e2cc269bbfb8bb49b5"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end