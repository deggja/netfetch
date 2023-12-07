class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters and provide a score"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/0.0.50/netfetch_0.0.50_darwin_amd64.tar.gz"
    sha256 "99e6f149065d5589d5306266320dbefc130c0a61f5ae2e67514de54df14f0d91"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end