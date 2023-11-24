class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/0.0.15/netfetch_0.0.15_darwin_amd64.tar.gz"
    sha256 "53bfbc3240986e4774ec2e7652be4bcac58c05e118d5443cbf0dd37ff8572722"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end