class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters and provide a score"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/0.0.40/netfetch_0.0.40_darwin_amd64.tar.gz"
    sha256 "8f82583cf1d35565b9af6cfb806467c0ddb04d1b8ed0184e148860673f30e54d"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end