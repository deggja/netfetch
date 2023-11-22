class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/0.0.1/netfetch_0.0.1_darwin_amd64.tar.gz"
    sha256 "0d9ba6bb2c0509b8a1d97bc680fccf48f2689e4f4195865aeb4dfa08356d8db0"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end
  