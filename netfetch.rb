class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters and provide a score"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/0.0.41/netfetch_0.0.41_darwin_amd64.tar.gz"
    sha256 "f95c1e8c0156b49578b40b78a3c05853e11e2a2a986e8f9b6ca851f6118e61f1"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end