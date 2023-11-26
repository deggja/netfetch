class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters and provide a score"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/0.0.37/netfetch_0.0.37_darwin_amd64.tar.gz"
    sha256 "ec6fbf079f8c313b0bbe562c68f9213d2f6a0d4630f3c8aa65dca852f92518c8"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end